/**
 * This template is a production ready boilerplate for developing with `PuppeteerCrawler`.
 * Use this to bootstrap your projects using the most up-to-date code.
 * If you're looking for examples or want to learn more, see README.
 */

const Apify = require('apify');

const { utils: { log } } = Apify;

Apify.main(async () => {
    const requestQueue = await Apify.openRequestQueue();
    await requestQueue.addRequest({ url: `https://www.zacks.com/logout.php` });

    // const proxyConfiguration = await Apify.createProxyConfiguration();

    const crawler = new Apify.PuppeteerCrawler({
        requestQueue,
        // proxyConfiguration,
        launchContext: {
            // Chrome with stealth should work for most websites.
            // If it doesn't, feel free to remove this.
            useChrome: true,
            stealth: true,
        },
        persistCookiesPerSession: true,
        useSessionPool: true,
        maxConcurrency: 1,
        handlePageTimeoutSecs: 120,
        // This function will be called for each URL to crawl.
        // Here you can write the Puppeteer scripts you are familiar with,
        // with the exception that browsers and pages are automatically managed by the Apify SDK.
        // The function accepts a single parameter, which is an object with the following fields:
        // - request: an instance of the Request class with information such as URL and HTTP method
        // - page: Puppeteer's Page object (see https://pptr.dev/#show=api-class-page)
        handlePageFunction: async ({ request, page }) => {

            await page.waitForSelector('#login input[name=username]');
            await page.type('#login input[name=username]', process.env.USER_ID);
            await page.type('#login input[name=password]', process.env.USER_PASS);
            await page.click('#login input[value=Login]');
            await page.waitForNavigation();

            await page.goto('https://www.zacks.com/screening/stock-screener');

            const iframe = await page.waitForSelector("#screenerContent");
            const frame = await iframe.contentFrame();

            await page.waitForTimeout(1000);

            await frame.waitForSelector('#my-screen-tab');
            await frame.click('#my-screen-tab');

            await page.waitForTimeout(1000);

            await frame.waitForSelector('#btn_run_137005');
            await frame.click('#btn_run_137005');

            await page.waitForTimeout(1000);

            await page._client.send('Page.setDownloadBehavior', { behavior: 'allow', downloadPath: './downloads' });

            await frame.waitForSelector('#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5');
            await frame.click('#screener_table_wrapper > div.dt-buttons > a.dt-button.buttons-csv.buttons-html5');

            await page.waitForTimeout(60000);
        },

        // This function is called if the page processing failed more than maxRequestRetries+1 times.
        handleFailedRequestFunction: async ({ request }) => {
            console.log(`Request ${request.url} failed too many times.`);
        },

        preNavigationHooks: [],
    });

    log.info('Starting the crawl.');
    await crawler.run();
    log.info('Crawl finished.');
});
