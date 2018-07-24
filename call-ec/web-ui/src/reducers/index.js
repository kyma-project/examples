const initialState = {
    isUserLoggedIn: false,
    user: {},
};

export default (state = initialState, action) => {
    switch (action.type) {
        case 'USER_LOGGED_IN':
            console.log(`Got USER_LOGGED_IN. User:  ${JSON.stringify(action.payload)}`);
            return {
                ...state,
                ...{
                    isUserLoggedIn: true,
                    user: action.payload,
                }
            };

        case 'USER_LOGGED_OUT':
            console.log('Got USER_LOGGED_OUT.');
            return {
                ...state,
                ...{
                    isUserLoggedIn: false,
                    user: {},
                }
            };
        default:
            return state;
    }
}