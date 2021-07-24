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
    currentCameraId : string;
    currentTransient : boolean;
}

export default class NodeEditPanel extends React.Component<NodeEditPanelProps, NodeEditPanelState> {

    constructor(props : NodeEditPanelProps) {
        super(props);
        let camId = "";
        if(props.node.CameraId != null) {
            camId = String(props.node.CameraId);
        }
        this.state = {
            currentCameraId: camId,
            currentTransient: props.node.IsTransient,
        };
    }

    private hasChanged() : boolean {
        let camId = "";
        if(this.props.node.CameraId != null) {
            camId = String(this.props.node.CameraId);
        }
        return camId != this.state.currentCameraId ||
                this.props.node.IsTransient != this.state.currentTransient;
    }

    render() {
        if(this.props.node) {

            let onCameraIdChange = (evt : React.ChangeEvent<HTMLInputElement>) => {
                let newVal = evt.currentTarget.value;
                this.setState((oldState : NodeEditPanelState) : NodeEditPanelState => {
                    return {
                      currentCameraId: newVal,
                      currentTransient: oldState.currentTransient,
                    };
                });
            };

            let onTransientChange = (evt : React.ChangeEvent<HTMLInputElement>) => {
                let newVal = evt.currentTarget.checked;
                this.setState((oldState : NodeEditPanelState) : NodeEditPanelState => {
                    return {
                        currentCameraId: oldState.currentCameraId,
                        currentTransient: newVal,
                    };
                });
            };

            let onSave = () => {
                let props : NodeProps = {
                    CameraId : this.state.currentCameraId == "" ?
                        null : Number(this.state.currentCameraId),
                    IsTransient: this.state.currentTransient,
                }
                this.props.ds.setNodeProps(this.props.node.mqttName, props)
            };

            return <div>
                <div className={styles.nodetitle}>Editing {this.props.node.name}</div>
                <div className={styles.nodepropsline}>
                    <span className={styles.nodepropstitle}>Camera Id: </span>
                    <span className={styles.nodepropsvalue}><input name="cameraId" value={this.state.currentCameraId} onChange={onCameraIdChange}/> </span>
                </div>


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
