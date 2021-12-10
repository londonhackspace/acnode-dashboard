import React from "react";
import ExtendedNodeRecord from "../extendednoderecord";

// yeah ok we -steal- borrow some styles from the detail panel
import styles from "./nodedetailspanel.module.css"
import DataSource from "../NodeDataSource";
import {NodeProps} from "../apiclient/dashapi";

interface NodeEditPanelProps {
    ds : DataSource
    node : ExtendedNodeRecord
}

interface NodeEditPanelState {
    currentTransient : boolean;
}

export default class NodeEditPanel extends React.Component<NodeEditPanelProps, NodeEditPanelState> {

    constructor(props : NodeEditPanelProps) {
        super(props);
        this.state = {
            currentTransient: props.node.IsTransient,
        };
    }

    private hasChanged() : boolean {
        return this.props.node.IsTransient != this.state.currentTransient;
    }

    render() {
        if(this.props.node) {

            let onTransientChange = (evt : React.ChangeEvent<HTMLInputElement>) => {
                let newVal = evt.currentTarget.checked;
                this.setState((oldState : NodeEditPanelState) : NodeEditPanelState => {
                    return {
                        currentTransient: newVal,
                    };
                });
            };

            let onSave = () => {
                let props : NodeProps = {
                    IsTransient: this.state.currentTransient,
                }
                this.props.ds.setNodeProps(this.props.node.mqttName, props)
            };

            return <div>
                <div className={styles.nodetitle}>Editing {this.props.node.name}</div>

                <div className={styles.nodepropsline}>
                    <span className={styles.nodepropstitle}>Is Transient: </span>
                    <span className={styles.nodepropsvalue}><input name="transient" type="checkbox" checked={this.state.currentTransient} onChange={onTransientChange}/> </span>
                </div>
                <button onClick={onSave} value="Save" disabled={!this.hasChanged()}>Save</button>
            </div>
        } else {
            return <div>Please select a node.</div>
        }
    }
}
