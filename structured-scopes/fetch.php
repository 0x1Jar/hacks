<?php
$url = "https://hackerone.com/graphql";
$authtoken = $argv[1]?? die('needs auth token');

$query = <<<QUERY
query Settings { 
    query{ 
        id,
        teams(first: 50 after: "%s") {
            pageInfo {
                hasNextPage,
                hasPreviousPage
            },
            edges{
                cursor,
                node{
                    _id,
                    handle,
                    structured_scopes {
                        edges {
                            node {
                                id,
                                asset_type,
                                asset_identifier,
                                eligible_for_submission,
                                eligible_for_bounty,
                                max_severity,
                                archived_at,
                                instruction
                            }
                        }
                    }
                }
            }
        }
    }   
}
QUERY;

$gen = function($cursor = "") use($query){
	return json_encode([
		'query' => sprintf($query, $cursor),
        'variables' => (object) []
	]);
};


$cursor = "";
$pageCount = 0;
$userAgent = "Mozilla/5.0 (compatible; StructuredScopesFetcher/1.0; +https://github.com/0x1Jar/new-hacks/tree/main/structured-scopes)"; // Example User-Agent

do {
    $pageCount++;
    fprintf(STDERR, "Fetching page %d (cursor: %s)\n", $pageCount, $cursor ?: "none");

    $params = [
        'http' => [
            'method' => 'POST',
            'header' => "Content-Type: application/json\r\n".
                        "Origin: https://hackerone.com\r\n".
                        "Referer: https://hackerone.com/programs\r\n".
                        "User-Agent: {$userAgent}\r\n".
                        "X-Auth-Token: {$authtoken}",
            'content' => $gen($cursor),
            'ignore_errors' => true // To handle HTTP errors manually
        ]
    ];
    $context = stream_context_create($params);
    $fp = @fopen($url, 'rb', false, $context); // Suppress fopen warnings, check manually

    if ($fp === false) {
        $error = error_get_last();
        die("Failed to open stream to {$url}. Error: " . ($error['message'] ?? 'Unknown error') . "\n");
    }

    $responseContent = stream_get_contents($fp);
    fclose($fp);

    if ($responseContent === false) {
        die("Failed to read stream content from {$url}.\n");
    }

    $result = json_decode($responseContent);

    // Check for JSON decoding errors
    if (json_last_error() !== JSON_ERROR_NONE) {
        die("Failed to decode JSON response. Error: " . json_last_error_msg() . "\nResponse snippet: " . substr($responseContent, 0, 200) . "\n");
    }

    // Check for GraphQL errors
    if (isset($result->errors) && is_array($result->errors) && count($result->errors) > 0) {
        $errorMessages = [];
        foreach ($result->errors as $gqlError) {
            $errorMessages[] = $gqlError->message ?? 'Unknown GraphQL error';
        }
        die("GraphQL API returned errors: " . implode("; ", $errorMessages) . "\n");
    }

    // Check for expected data structure
    if (!isset($result->data->query->teams->pageInfo) || !isset($result->data->query->teams->edges)) {
        die("Unexpected API response structure. Response snippet: " . substr($responseContent, 0, 500) . "\n");
    }
    
    $hasNextPage = $result->data->query->teams->pageInfo->hasNextPage;

    foreach ($result->data->query->teams->edges as $edge){
        $cursor = $edge->cursor;
        // Ensure node and structured_scopes exist before iterating
        if (!isset($edge->node->structured_scopes->edges)) {
            fprintf(STDERR, "Warning: Team %s (handle: %s) missing structured_scopes->edges. Skipping.\n", $edge->node->_id ?? 'N/A', $edge->node->handle ?? 'N/A');
            continue;
        }
        foreach ($edge->node->structured_scopes->edges as $scope){
            $scope = $scope->node;
            if (!$scope->eligible_for_submission){
                continue;
            }
            if (strToLower($scope->asset_type) != "url"){
                continue;
            }

            echo $scope->asset_identifier.PHP_EOL;
        } 
    }

} while($hasNextPage);
