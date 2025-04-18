import { collectJSDelivrStats } from './scrapers/jsdelivrScraper';
import { collectNPMStats } from './scrapers/npmScraper';

async function main() {
    try {
        console.log('Starting NPM stats collection...');
        await collectNPMStats();
        console.log('NPM stats collection completed successfully');

        console.log('Starting JSDelivr stats collection...');
        await collectJSDelivrStats();
        console.log('JSDelivr stats collection completed successfully');
    } catch (error) {
        console.error('Error collecting stats:', error);
        process.exit(1);
    }
}

main();