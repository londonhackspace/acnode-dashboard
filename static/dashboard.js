const urlPrefix = "/api/";

let nodes = [];

function getApiData(endpoint) {
    const url = urlPrefix + endpoint;

    let p = new Promise(function (resolve, reject) {
        let xhttp = new XMLHttpRequest();
        xhttp.open("GET", url, true);
        xhttp.onreadystatechange = function() {
            if(this.readyState == 4) {
                if(this.status == 200) {
                    resolve(JSON.parse(this.responseText));
                } else {
                    reject(this.status);
                }
            }
        }

        xhttp.send();
    });

    return p;
}

function getNodes() {
    return getApiData("nodes");
}

function getNode(name) {
    return getApiData("nodes/" + name);
}

function makeNodeClassName(node) {
    if(node.LastSeen == -1) {
        return "node, node-unknown";
    }
    if(node.LastSeen < 180) {
        return "node node-good";
    }

    return "node node-bad";
}

function makeNodeHtml(node) {
    let root = document.createElement("div");

    let name = document.createElement("div");
    name.innerText = "Node: " + node.name;

    root.append(name);
    root.append(document.createElement("br"));

    let lastSeen = document.createElement("div")
    if(node.LastSeen == -1) {
        lastSeen.innerText = "Last seen: Never";
    } else {
        lastSeen.innerText = "Last seen: " + node.LastSeen + " seconds ago";
    }

    root.append(lastSeen);
    root.append(document.createElement("br"));

    if(node.Status != "") {
        let status = document.createElement("div");
        status.innerText = "Status: " + node.Status;
        root.append(status);
    }

    if(node.Version != "") {
        let version = document.createElement("div");
        let link = document.createElement("a");
        let ver = node.Version;
        version.append(document.createTextNode("Version: "))
        let vertext = "";
        if(ver.includes("-")) {
            let split = ver.split("-");
            ver = split[0];
            vertext = "-"+split[1];
        }
        link.href = "https://github.com/londonhackspace/acnode-cl/commits/" + ver;
        link.innerText = ver;
        version.append(link);
        version.append(document.createTextNode(vertext))
        root.append(version);
    }

    root.append(document.createElement("br"));

    if((node.MemFree + node.MemUsed) > 0) {
        let chartarea = document.createElement("canvas");
        chartarea.className = "nodechart";
        let chart = new Chart(chartarea, {
            type: 'doughnut',
            data: {
                labels: ["Memory Used", "Memory Free"],
                datasets: [
                    {
                        data: [node.MemUsed, node.MemFree],
                        backgroundColor: [
                            'rgb(255, 99, 132)',
                            'rgb(54, 162, 235)'
                        ],
                        hoverOffset: 4
                    }
                ]
            }
        });

        root.append(chartarea);
    }

    return root;
}

function addNode(nodeName) {
    nodes.push(nodeName);
    let nc = document.getElementById("nodecontainer");
    let el = document.createElement("div")
    el.id = "node-" + nodeName + "-container";
    el.className = "node node-unknown";

    let titleLine = document.createElement("div");
    titleLine.className = "nodetitle";
    titleLine.innerText = "Node: " + nodeName;
    el.append(titleLine);

    nc.append(el);
}

function updateNodes() {
    getNodes().then((res) => {
        for(node of res) {
            if(!nodes.includes(node)) {
                addNode(node);
            }
        }
    }).then(() => {
        let fulldata = [];
        for(node of nodes) {
            fulldata.push(getNode(node));
        }
        return Promise.all(fulldata);
    }).then((nodeData) => {
        for(node of nodeData) {
            let el = document.getElementById("node-" + node.mqttName + "-container");
            el.innerHTML = "";
            el.append(makeNodeHtml(node));
            el.className = makeNodeClassName(node);
        }
    });
}

window.onload = updateNodes;

setTimeout(updateNodes, 10000);