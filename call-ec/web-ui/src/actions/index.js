
export const logIn = (userData) => {
    return {
        type: 'USER_LOGGED_IN',
        payload: userData
    }
};

export const logOut = () => {
    return {
        type: 'USER_LOGGED_OUT'
    }
};