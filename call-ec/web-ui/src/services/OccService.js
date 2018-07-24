const esc = encodeURIComponent;

class OccService {

    /**
     *
     * @param options.userId identifier of the user
     * @param options.threshold threshold for number of orders required to have promotion granted
     * @param options.ecIdToken ID token issued by EC for the user
     * @param options.ecAccessToken access token issued by EC
     * @returns {Promise<any>}
     */
    getPromotion = async (options) => {

        const getPromotionUrl = process.env.REACT_APP_CALCULATE_PROMOTION_URL;

        const params = {
            'user-id': options.userId,
            threshold: options.threshold
        };
        const query = Object.keys(params)
            .map(k => esc(k) + '=' + esc(params[k]))
            .join('&');
        const url = `${getPromotionUrl}?${query}`;

        const requestOptions = {
            method: 'POST',
            headers: {
                "Content-type": "application/x-www-form-urlencoded",
                "Authorization": `Bearer ${options.ecIdToken}`,
                "occ-token": options.ecAccessToken,
            },
            mode: 'cors'
        };

        const response = await fetch(url, requestOptions);
        return await response.json();
    }
}

export default OccService;