import axios from 'axios';
import { CSVWriter } from '../../utils/csvWriter';

interface VersionDownloads {
    [version: string]: number;
}

interface NPMResponse {
    downloads: VersionDownloads;
}

async function fetchNPMStats(packageName: string): Promise<NPMResponse> {
    const url = `https://api.npmjs.org/versions/@astrouxds%2F${packageName}/last-week`;

    try {
        const response = await axios.get<NPMResponse>(url);
        return response.data;
    } catch (error) {
        if (axios.isAxiosError(error)) {
            if (error.response?.status) {
                throw new Error(`API returned non-200 status: ${error.response.status}`);
            }
            throw new Error(`Error making request: ${error.message}`);
        }
        throw error;
    }
}

export async function collectNPMStats(): Promise<void> {
    const packages = ['astro-web-components', 'react'];
    const csvWriter = new CSVWriter('npm_stats.csv');
    const currentDate = new Date().toISOString().split('T')[0];

    for (const pkg of packages) {
        try {
            const stats = await fetchNPMStats(pkg);
            
            for (const [version, downloads] of Object.entries(stats.downloads)) {
                await csvWriter.appendData(currentDate, pkg, version, downloads);
            }
        } catch (error) {
            console.error(`Error fetching stats for ${pkg}:`, error);
            continue;
        }
    }
}
