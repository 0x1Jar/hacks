const puppeteer = require('puppeteer');

puppeteer.launch({ignoreHTTPSErrors: true}).then(async browser => {
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    await page.setUserAgent('Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36');

    page.on('request', interceptedRequest => {
        interceptedRequest.continue().catch(err => {
            // console.error('Error continuing request:', err.message); // Optional: log if continue fails
        });
    });

    page.on('response', async resp => {
        const headers = resp.headers();
        const contentType = headers['content-type'];

        if (!contentType) {
            return;
        }

        if (contentType.match(/(javascript|json)/i)) {
            console.log(resp.url());
        }
    });

    let url = process.argv[2];
    if (!url) {
        console.error("Please provide a URL as a command-line argument.");
        await browser.close();
        process.exit(1);
    }

    try {
        await page.goto(url, { waitUntil: 'networkidle2', timeout: 30000 }); // Added waitUntil and timeout
    } catch (e) {
        console.error(`Error navigating to ${url}: ${e.message}`);
    } finally {
        await browser.close();
    }
}).catch(err => {
    console.error("Puppeteer launch failed:", err.message);
    process.exit(1);
});
