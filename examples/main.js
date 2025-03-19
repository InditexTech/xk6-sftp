import xk6sftp from "k6/x/sftp";
import { sleep, group, check } from "k6";

const client = xk6sftp.newClient("user", "pwd", "localhost", 3322);
const arrayFiles = ["utf-8.txt", "utf-16.txt", "binary.png"];

export default function () {
    const remotePath = "/wwwroot";
    // Local path is relative to the repository root path
    const localPath = "./examples/test-data";

    arrayFiles.forEach(file => {
        group(`Upload ${file}`, function () {
            const r = client.uploadFile(`${localPath}/${file}`, `${remotePath}/${file}`);
            console.log(`Upload file ${file} - result: ${JSON.stringify(r)}`);
            check(r, {
                "is success": (r) => r.success === true,
            });
        });
        sleep(0.3);
    });

    arrayFiles.forEach(file => {
        group(`Download ${file}`, function () {
            const r = client.downloadFile(`${remotePath}/${file}`, `${localPath}/downloaded/${file}`);
            console.log(`Download file ${file} - result: ${JSON.stringify(r)}`);
            check(r, {
                "is success": (r) => r.success === true,
            });
        });
        sleep(0.3);
    });

    arrayFiles.forEach(file => {
        group(`Delete ${file}`, function () {
            const r = client.deleteFile(`${remotePath}/${file}`);
            console.log(`Delete file ${file} - result: ${JSON.stringify(r)}`);
            check(r, {
                "is success": (r) => r.success === true,
            });
        });
        sleep(0.3);
    })
}

export function teardown() {
    if (client !== null) {
        client.close();
    }
}

