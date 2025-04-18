import { createObjectCsvWriter } from 'csv-writer';

export class CSVWriter {
    private writer;

    constructor(filename: string) {
        this.writer = createObjectCsvWriter({
            path: filename,
            header: [
                { id: 'date', title: 'Date' },
                { id: 'package', title: 'Package' },
                { id: 'version', title: 'Version' },
                { id: 'downloads', title: 'Downloads' }
            ]
        });
    }

    async appendData(date: string, packageName: string, version: string, downloads: number) {
        return this.writer.writeRecords([{
            date,
            package: packageName,
            version,
            downloads
        }]);
    }
}
