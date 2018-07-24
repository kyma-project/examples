import React, {Component, Fragment} from 'react';
import {withStyles} from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import TextField from '@material-ui/core/TextField';
import OccService from "../services/OccService";

const styles = theme => ({
    root: {
        width: '100%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(15),
        fontWeight: theme.typography.fontWeightRegular,
    },
    errorMsg: {
        color: "red",
        fontWeight: theme.typography.fontWeightMedium
    },
    textField: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
        width: 200,
    },
    promotionInfo: {
        fontWeight: theme.typography.fontWeightMedium
    },
});

class Lambda extends Component {

    constructor(props) {

        super(props);

        this.user = props.withUser;

        this.state = {
            isLoading: false,
            result: null,
            error: null,
            threshold: 10
        };

        this.occService = new OccService();
    }

    render() {

        const {classes} = this.props;
        const {isLoading, error, result} = this.state;

        return (
            <Fragment>
                <h1>Get Promotions From EC</h1>

                <ExpansionPanel>
                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon/>}>
                        <Typography className={classes.heading}>Authentication</Typography>
                    </ExpansionPanelSummary>
                    <ExpansionPanelDetails className={classes.details}>
                        <ul>
                            <li><b>User:</b> {this.user.id}</li>
                            <li><b>ID token:</b> {this.user.idToken}</li>
                            <li><b>Access token:</b> {this.user.accessToken}</li>
                        </ul>
                    </ExpansionPanelDetails>
                </ExpansionPanel>


                <ExpansionPanel>
                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon/>}>
                        <Typography className={classes.heading}>The Function</Typography>
                    </ExpansionPanelSummary>
                    <ExpansionPanelDetails className={classes.details}>
                        <div>-
                            <TextField
                                id="threshold"
                                label="Threshold"
                                value={this.state.threshold}
                                onChange={this.handleChange('threshold')}
                                type="number"
                                className={classes.textField}
                                InputLabelProps={{
                                    shrink: true,
                                }}
                                margin="normal"
                            />
                        </div>
                        <div>
                            {(!error && !result && !isLoading) && (<Typography>No data</Typography>)}
                            {isLoading && (<Typography>Loading...</Typography>)}
                            {error && (<Typography className={classes.errorMsg}>Error: {error}</Typography>)}
                            {result && (<Typography>Response: {JSON.stringify(result)}</Typography>)}
                            {(result && result.promotion) && (
                                <Typography className={classes.promotionInfo}>GOT PROMOTION!</Typography>)}
                            {(result && !result.promotion) && (
                                <Typography className={classes.promotionInfo}>Promotion not granted :( Try
                                    again!</Typography>)}
                        </div>
                    </ExpansionPanelDetails>
                    <Divider/>
                    <ExpansionPanelActions>
                        <Button onClick={this.callLambda}
                                disabled={isLoading}
                                variant="raised" color="primary"
                                title="Call the function">Call the function</Button>
                    </ExpansionPanelActions>
                </ExpansionPanel>
            </Fragment>
        );
    }

    handleChange = name => event => {
        this.setState({
            [name]: event.target.value,
        });
    };

    callLambda = async () => {

        this.setState({
            isLoading: true,
            result: null,
            error: null
        });

        const options = {
            userId: this.user.id,
            threshold: this.state.threshold,
            ecIdToken: this.user.idToken,
            ecAccessToken: this.user.accessToken
        };

        try {
            const result = await this.occService.getPromotion(options);
            this.setState({
                isLoading: false,
                result: result
            });
        }
        catch(err) {
            console.log(err);
            this.setState({
                isLoading: false,
                error: err.toString()
            });
        }
    }
}

export default withStyles(styles)(Lambda);