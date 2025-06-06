const puppeteer = require('puppeteer');
const readline = require('readline');
const fs = require('fs');
const yargs = require('yargs/yargs');
const { hideBin } = require('yargs/helpers');

const DEFAULT_CONCURRENCY = 10;

async function checkRedirect(browser, url) {
    const page = await browser.newPage();
    try {
        await page.setExtraHTTPHeaders({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        });
        await page.goto(url, { waitUntil: 'networkidle2', timeout: 30000 }); // Added timeout

        const destination = await page.evaluate(() => {
            return { "domain": document.domain, "href": document.location.href };
        });

        const originalUrl = new URL(url);

        if (originalUrl.host.toLowerCase() !== destination.domain.toLowerCase() || originalUrl.pathname !== new URL(destination.href).pathname) {
            console.log(`${url} redirects to ${destination.href}`);
        } else {
            console.log(`${url} does not redirect`);
        }
    } catch (e) {
        console.error(`Error checking ${url}: ${e.message}`);
    } finally {
        await page.close();
    }
}

async function main() {
    const argv = yargs(hideBin(process.argv))
        .option('c', {
            alias: 'concurrency',
            type: 'number',
            description: 'Number of concurrent pages',
            default: DEFAULT_CONCURRENCY
        })
        .option('f', {
            alias: 'file',
            type: 'string',
            description: 'Input file with URLs (one per line)'
        })
        .help()
        .alias('help', 'h')
        .argv;

    const concurrency = argv.concurrency > 0 ? argv.concurrency : DEFAULT_CONCURRENCY;
    const inputFile = argv.file;

    let urls = [];
    const
    inputStream = inputFile ? fs.createReadStream(inputFile) : process.stdin;

    const rl = readline.createInterface({
        input: inputStream,
        crlfDelay: Infinity
    });

    for await (const line of rl) {
        if (line.trim()) {
            urls.push(line.trim());
        }
    }

    if (urls.length === 0) {
        console.log("No URLs to process.");
        return;
    }

    console.log(`Processing ${urls.length} URLs with concurrency ${concurrency}...`);

    const browser = await puppeteer.launch({ 
        ignoreHTTPSErrors: true,
        headless: true, // Explicitly set, can be 'new' for newer versions
        args: ['--no-sandbox', '--disable-setuid-sandbox'] // Common args for CI environments
    });
    
    const queue = [...urls];
    let activePromises = 0;
    const results = []; // To store results or handle them as they come

    const processNext = async () => {
        if (queue.length === 0) {
            return; // All URLs processed
        }
        if (activePromises >= concurrency) {
            return; // Wait for a slot
        }

        const url = queue.shift();
        activePromises++;

        try {
            await checkRedirect(browser, url);
        } catch (e) {
            // Error already logged in checkRedirect, or log here if preferred
        } finally {
            activePromises--;
            // If queue is empty and no active promises, or if more items to process, call processNext
            if (queue.length > 0 || activePromises > 0) {
                 processNext(); // Immediately try to process next if slots/items available
            }
        }
    };

    // Start initial batch of workers
    for (let i = 0; i < concurrency && queue.length > 0; i++) {
        processNext();
    }

    // Keep the main function alive until all URLs are processed
    // This can be managed by checking if queue is empty and activePromises is 0
    // A more robust way is to use a Promise.all for all tasks if they return promises.
    // For now, we'll rely on the activePromises and queue length.

    // A simple polling mechanism to wait for completion
    await new Promise(resolve => {
        const interval = setInterval(() => {
            if (queue.length === 0 && activePromises === 0) {
                clearInterval(interval);
                resolve();
            }
        }, 100); // Check every 100ms
    });

    console.log("All URLs processed.");
    await browser.close();
}

main().catch(err => {
    console.error("Unhandled error in main:", err);
    process.exit(1);
});
