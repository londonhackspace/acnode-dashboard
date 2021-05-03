
export function restGetRequest(url:string) : Promise<string> {
    let p = new Promise<string>((resolve, reject) => {
        let httpReq = new XMLHttpRequest();
        httpReq.open("GET", url, true)
        httpReq.onreadystatechange = function() {
            if(this.readyState == 4) {
                if(this.status >= 200 && this.status < 300) {
                    resolve(this.response);
                } else {
                    reject(this.status);
                }
            }
        }
        httpReq.send();
    });
    return p;
}