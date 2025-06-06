package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/machinebox/graphql"
)

// Corrected structs for GraphQL response with pagination
type graphQLResponse struct {
	Query struct {
		Program programPayload `json:"Program"` // Matches the alias "Program" in the GraphQL query
	} `json:"query"`
}

type programPayload struct {
	Handle     string     `json:"handle"`
	InScope    scopeAssets `json:"InScope"`    // Matches alias
	OutOfScope scopeAssets `json:"OutOfScope"` // Matches alias
}

type scopeAssets struct {
	Edges    []assetEdge `json:"Assets"` // Matches alias "Assets" which contains edges
	PageInfo pageInfo    `json:"pageInfo"`
}

type assetEdge struct {
	Node   asset  `json:"Asset"` // Matches alias "Asset" which is the node
	Cursor string `json:"cursor"`
}

type asset struct {
	Identifier string `json:"asset_identifier"`
	Type       string `json:"asset_type"`
	// Potentially other fields like rendered_instruction, max_severity, eligible_for_bounty if needed from node
}

type pageInfo struct {
	HasNextPage     bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

func (a asset) Domain() (string, error) {
	if a.Type != "URL" {
		return "", fmt.Errorf("asset with identifier %s is not a URL", a.Identifier)
	}

	i := strings.ToLower(a.Identifier)

	// If it has a scheme then parse it and return the hostname
	if a.hasScheme() {
		u, err := url.Parse(i)
		if err == nil {
			return u.Hostname(), nil
		}
	}

	if a.isWildcard() {
		return strings.TrimLeft(i, "*.%"), nil
	}

	return i, nil
}

func (a asset) isWildcard() bool {
	if len(a.Identifier) < 2 {
		return false
	}

	if a.Identifier[0] == '*' {
		return true
	}

	if a.Identifier[0] == '.' {
		return true
	}

	if a.Identifier[0] == '%' {
		return true
	}

	return false
}

func (a asset) hasScheme() bool {
	i := strings.ToLower(a.Identifier)
	if len(i) < 6 {
		return false
	}

	if i[:5] == "http:" {
		return true
	}

	if len(i) < 7 {
		return false
	}

	if i[:6] == "https:" {
		return true
	}

	return false
}

func main() {

	var risky bool
	flag.BoolVar(&risky, "risky", false, "treat all domains as wildcards")

	var appendScope bool
	flag.BoolVar(&appendScope, "append-scope", false, "append to the scope file instead of replacing it")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Error: Program handle is required.")
		fmt.Printf("Usage: %s [options] <program_handle>\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	programHandle := flag.Arg(0)

	graphQLToken := os.Getenv("H1_GRAPHQL_TOKEN")
	if graphQLToken == "" {
		fmt.Println("H1_GRAPHQL_TOKEN not set. Go to https://hackerone.com/current_user/graphql_token.json to get one")
		return
	}

	fmt.Printf("Fetching scope for program: %s\n", programHandle)

	client := graphql.NewClient("https://hackerone.com/graphql")
	ctx := context.Background()

	var allInScopeAssets []assetEdge
	var allOutOfScopeAssets []assetEdge

	// Fetch In-Scope Assets with Pagination
	var inScopeCursor *string
	fmt.Println("Fetching in-scope assets...")
	for {
		req := buildGraphQLRequest(programHandle, graphQLToken, inScopeCursor, nil)
		var pageResp graphQLResponse // Use the new top-level response struct
		if err := client.Run(ctx, req, &pageResp); err != nil {
			log.Fatalf("Error fetching in-scope assets: %v", err)
		}
		allInScopeAssets = append(allInScopeAssets, pageResp.Query.Program.InScope.Edges...)
		if !pageResp.Query.Program.InScope.PageInfo.HasNextPage || len(pageResp.Query.Program.InScope.Edges) == 0 {
			break
		}
		lastInScopeAssetEdge := pageResp.Query.Program.InScope.Edges[len(pageResp.Query.Program.InScope.Edges)-1]
		inScopeCursor = &lastInScopeAssetEdge.Cursor
	}
	fmt.Printf("Fetched %d in-scope assets.\n", len(allInScopeAssets))

	// Fetch Out-Of-Scope Assets with Pagination
	var outOfScopeCursor *string
	fmt.Println("Fetching out-of-scope assets...")
	for {
		req := buildGraphQLRequest(programHandle, graphQLToken, nil, outOfScopeCursor)
		var pageResp graphQLResponse // Use the new top-level response struct
		if err := client.Run(ctx, req, &pageResp); err != nil {
			log.Fatalf("Error fetching out-of-scope assets: %v", err)
		}
		allOutOfScopeAssets = append(allOutOfScopeAssets, pageResp.Query.Program.OutOfScope.Edges...)
		if !pageResp.Query.Program.OutOfScope.PageInfo.HasNextPage || len(pageResp.Query.Program.OutOfScope.Edges) == 0 {
			break
		}
		lastOutOfScopeAssetEdge := pageResp.Query.Program.OutOfScope.Edges[len(pageResp.Query.Program.OutOfScope.Edges)-1]
		outOfScopeCursor = &lastOutOfScopeAssetEdge.Cursor
	}
	fmt.Printf("Fetched %d out-of-scope assets.\n", len(allOutOfScopeAssets))

	// truncate the scopes file by default
	sf, df, wf, err := openOutputFiles(appendScope)
	if err != nil {
		log.Fatalf("Error opening output files: %v", err)
	}
	defer sf.Close()
	defer df.Close()
	defer wf.Close()

	processedDomains := make(map[string]bool) // To avoid duplicate entries in domains/wildcards files

	fmt.Println("Processing in-scope assets...")
	for _, edge := range allInScopeAssets { // Iterate over edges
		a := edge.Node // Get asset from node
		d, err := a.Domain()
		if err != nil {
			fmt.Printf("Info: Skipping in-scope asset '%s' (type: %s): %v\n", a.Identifier, a.Type, err)
			continue
		}

		isWildcardAsset := a.isWildcard() || risky
		if isWildcardAsset {
			fmt.Fprintf(sf, ".*\\.%s$\n", d) // Burp scope regex for wildcard
			if !processedDomains[d+"_wildcard"] { // Check processedDomains before writing
				fmt.Fprintf(wf, "%s\n", d)
				processedDomains[d+"_wildcard"] = true
			}
		}

		// Always add exact match for the domain itself to .scope and domains
		fmt.Fprintf(sf, "^%s$\n", d) // Burp scope regex for exact domain
		if !processedDomains[d] { // Check processedDomains before writing
			fmt.Fprintf(df, "%s\n", d)
			processedDomains[d] = true
		}
	}

	fmt.Println("Processing out-of-scope assets...")
	for _, edge := range allOutOfScopeAssets { // Iterate over edges
		a := edge.Node // Get asset from node
		d, err := a.Domain()
		if err != nil {
			fmt.Printf("Info: Skipping out-of-scope asset '%s' (type: %s): %v\n", a.Identifier, a.Type, err)
			continue
		}

		if a.isWildcard() {
			fmt.Fprintf(sf, "!.*\\.%s$\n", d) // Burp scope regex for out-of-scope wildcard
		}
		// Always add exact match for out-of-scope assets
		fmt.Fprintf(sf, "!^%s$\n", d) // Burp scope regex for out-of-scope exact domain
	}
	fmt.Println("Scope files generated successfully.")
}

// Helper function to build GraphQL request
func buildGraphQLRequest(handle string, token string, inScopeAfter *string, outOfScopeAfter *string) *graphql.Request {
	// Note: The GraphQL query is simplified here for brevity in the diff.
	// The actual query needs to be structured to accept $inScopeAfterCursor and $outOfScopeAfterCursor.
	// For a full implementation, the query string itself would need modification.
	// This is a placeholder to illustrate the pagination call structure.
	// The original query string needs to be parameterized for 'after' cursors.
	// For example, InScope:structured_scopes(first:$first_0, after:$inScopeAfterCursor, ...)
	// And OutOfScope:structured_scopes(first:$first_0, after:$outOfScopeAfterCursor, ...)

	scopesQuery := `
		query Team_assets($first_0:Int! $handle:String! $inScopeAfterCursor:String $outOfScopeAfterCursor:String) {
			query {
				id,
				...F0
			}
		}
		fragment F0 on Query {
			me {
				Membership:membership(team_handle:$handle) {
					permissions,
					id
				},
				id
			},
			Program:team(handle:$handle) {
				handle,
				_structured_scope_versions2ZWKHQ:structured_scope_versions(archived:false) {
					max_updated_at
				},
				InScope:structured_scopes(first:$first_0, after:$inScopeAfterCursor, archived:false,eligible_for_submission:true) {
					Assets:edges {
						Asset:node {
							id,
							asset_type,
							asset_identifier,
							rendered_instruction,
							max_severity,
							eligible_for_bounty
						},
						cursor # Ensure cursor is part of the AssetNode in the query
					},
					pageInfo {
						hasNextPage,
						hasPreviousPage
					}
				},
				OutOfScope:structured_scopes(first:$first_0, after:$outOfScopeAfterCursor, archived:false,eligible_for_submission:false) {
					Assets:edges {
						Asset:node {
							id,
							asset_type,
							asset_identifier,
							rendered_instruction
						},
						cursor # Ensure cursor is part of the AssetNode in the query
					},
					pageInfo {
						hasNextPage,
						hasPreviousPage
					}
				},
				id
			},
		id
		}
	`
	req := graphql.NewRequest(scopesQuery)
	req.Var("first_0", 50) // Reduced page size for more efficient pagination testing, can be 250
	req.Var("handle", handle)
	if inScopeAfter != nil {
		req.Var("inScopeAfterCursor", *inScopeAfter)
	} else {
		req.Var("inScopeAfterCursor", nil) // Explicitly set to null if not provided
	}
	if outOfScopeAfter != nil {
		req.Var("outOfScopeAfterCursor", *outOfScopeAfter)
	} else {
		req.Var("outOfScopeAfterCursor", nil) // Explicitly set to null if not provided
	}
	req.Header.Set("X-Auth-Token", token)
	return req
}

// Helper function to open output files
func openOutputFiles(appendScope bool) (sf, df, wf *os.File, err error) {
	scopeFlags := os.O_TRUNC | os.O_CREATE | os.O_WRONLY
	if appendScope {
		scopeFlags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}

	sf, err = os.OpenFile(".scope", scopeFlags, 0644)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to open .scope file: %w", err)
	}

	df, err = os.OpenFile("domains", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		sf.Close() // Close already opened file on error
		return nil, nil, nil, fmt.Errorf("failed to open domains file: %w", err)
	}

	wf, err = os.OpenFile("wildcards", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		sf.Close()
		df.Close() // Close already opened files on error
		return nil, nil, nil, fmt.Errorf("failed to open wildcards file: %w", err)
	}
	return sf, df, wf, nil
}
