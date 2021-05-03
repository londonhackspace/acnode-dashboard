import { restGetRequest} from "./utils";

interface CommitUser {
    name: string
    email: string
    date: string
}

interface CommitObject {
    url: string
    author: CommitUser
    committer: CommitUser
    message: string
}

interface CommitInfo {
    url : string
    sha: string
    html_url: string
    comments_url : string
    commit : CommitObject
}

// Extract only the fields we care about
function trimCommitInfo(ci : CommitInfo) : CommitInfo {
    return {
        url: ci.url,
        sha: ci.sha,
        comments_url: ci.comments_url,
        html_url: ci.html_url,
        commit: {
            url: ci.commit.url,
            author: ci.commit.author,
            committer: ci.commit.committer,
            message: ci.commit.message,
        }
    };
}

class GitCommitInfo {
    private cache : Map<string,CommitInfo> = new Map();

    private previous : Promise<CommitInfo> =
        new Promise((resolve) => resolve(null));

    private getCommitInfo(hash : string) : Promise<CommitInfo>{
        console.log("Looking up commit " + hash)
        return restGetRequest("https://api.github.com/repos/londonhackspace/acnode-cl/commits/"+hash)
            .then((data) => {
               return JSON.parse(data) ;
            }, function (errcode) {
                console.log("Error getting commit - " + errcode);
                return new Promise<CommitInfo>((resolve,reject) => reject(errcode));
            });
    }

    getCommit(hash : string ) : Promise<CommitInfo> {
        if(hash.endsWith("-dirty")) {
            hash = hash.substr(0, hash.length-"-dirty".length);
        }

        if(this.cache.has(hash)) {
            return new Promise<CommitInfo>((resolve => resolve(this.cache.get(hash))));
        }

        // sequentially process requests so the cache is as effective as it can be
        let p = this.previous.then(() => {
            // recheck cache
            if(this.cache.has(hash)) {
                return this.cache.get(hash);
            }
            return this.getCommitInfo(hash)
        }, () => null).then((result) => {
            let trimmed = result ? trimCommitInfo(result) : null;
            this.cache.set(hash, trimmed);
            return trimmed;
        }, () => {
            this.cache.set(hash, null);
            return null;
        });
        this.previous = p;
        return p;
    }
}
let GithubCommitInfo = new GitCommitInfo();
export default GithubCommitInfo;