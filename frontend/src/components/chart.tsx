import React from "react";
import ReactDOM from "react-dom";
import {Chart as ChartJS,
    DoughnutController,
    ArcElement,
    ChartTypeRegistry} from "chart.js"
import {MainFrameProps} from "../mainframe/main";

ChartJS.register(DoughnutController, ArcElement);

export interface ChartProps {
    type : keyof ChartTypeRegistry;
    data : Map<string, number>;
    colors: string[];
}

interface ChartState {

}

export default class Chart extends React.Component<ChartProps, ChartState> {
    private chart : ChartJS;

    private getDataKeys() : string[] {
        let res : string[] = [];
        for(let key of this.props.data.keys()) {
            res.push(key);
        }
        return res;
    }

    private getDataValues() : number[] {
        let res : number[] = [];
        for(let val of this.props.data.values()) {
            res.push(val);
        }
        return res;
    }

    private update() {
        console.log("Updating chart");
        this.chart.data.datasets[0].data = this.getDataValues();
        this.chart.data.labels = this.getDataKeys();
        this.chart.data.datasets[0].backgroundColor = this.props.colors;
        this.chart.update();
    }

    componentDidUpdate(prevProps : Readonly<ChartProps>) {
        this.update();
    }

    componentDidMount() {
        let node : HTMLCanvasElement = ReactDOM.findDOMNode(this) as HTMLCanvasElement;
        if(this.chart == null) {
            this.chart = new ChartJS(node.getContext('2d'), {
                type: this.props.type,
                data: {
                    labels: this.getDataKeys(),
                    datasets: [
                        {
                            data: this.getDataValues(),
                            backgroundColor: this.props.colors,
                            hoverOffset: 4
                        }
                    ]
                }
            });
        } else {
            this.update();
        }
    }

    componentWillUnmount() {
        this.chart = null;
    }

    render() {
        return <canvas></canvas>;
    }
}