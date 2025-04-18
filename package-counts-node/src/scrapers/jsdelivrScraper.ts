import axios from 'axios';
import { CSVWriter } from '../../utils/csvWriter';

export interface StatsResponse {
    rank: number;
    typeRank: number;
    total: number;
    versions: {
        [key: string]: VersionData;
    };
}

export interface VersionData {
    total: number;
    dates: {
        [key: string]: number;
    };
}

export interface PackageStats {
    name: string;
    version: string;
    downloads: {
        [key: string]: number;
    };
}

export async function fetchJSDelivrStats(packageName: string): Promise<PackageStats[]> {
    const url = `https://data.jsdelivr.com/v1/package/npm/${packageName}/stats`;
    
    try {
        const response = await axios.get<StatsResponse>(url);
        const allVersionStats: PackageStats[] = [];

        // Iterate through all versions
        Object.entries(response.data.versions).forEach(([version, data]) => {
            allVersionStats.push({
                name: packageName,
                version,
                downloads: data.dates
            });
        });

        return allVersionStats;
    } catch (error) {
        if (axios.isAxiosError(error)) {
            throw new Error(`Error fetching stats: ${error.message}`);
        }
        throw error;
    }
}

export async function collectJSDelivrStats() {
    const packages = ['@astrouxds/react', '@astrouxds/astro-web-components'];
    const csvWriter = new CSVWriter('jsdelivr_stats.csv');

    for (const pkg of packages) {
        try {
            const stats = await fetchJSDelivrStats(pkg);
            
            for (const stat of stats) {
                for (const [date, downloads] of Object.entries(stat.downloads)) {
                    await csvWriter.appendData(date, stat.name, stat.version, downloads);
                }
            }
        } catch (error) {
            console.error(`Error fetching stats for ${pkg}:`, error);
        }
    }
}

// Run the script
// collectJSDelivrStats()
//     .then(() => console.log('Done collecting stats'))
//     .catch(error => console.error('Error:', error));
