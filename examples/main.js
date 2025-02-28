import xk6sftp from "k6/x/sftp";
import { sleep, group, check } from "k6";

const client = xk6sftp.newClient("user", "pwd", "localhost", 3322);
const arrayFiles = ["utf-8.txt", "utf-16.txt", "binary.png"];

export default function () {
    const remotePath = "/wwwroot";

    // Relative path from parent director, aka, make run
    arrayFiles.forEach(file => {
        group(`Upload ${file}`, function () {
            const r = client.uploadFile(`./examples/test-data/${file}`, `${remotePath}/${file}`);
            //console.log(r);
            check(r, {
                "is success": (r) => r.success === true,
            });
        });
        sleep(0.3);
    });

    arrayFiles.forEach(file => {
        group(`Download ${file}`, function () {
            const r = client.downloadFile(`${remotePath}/${file}`, `./examples/test-data/downloaded/${file}`);
            //console.log(r);
            check(r, {
                "is success": (r) => r.success === true,
            });
        });

        sleep(0.3);
    });
}

export function teardown() {
    if (client !== null) {
        client.close();
    }
}

