const request = require('request');

const USER_ID_PARAM = "user-id";
const ORDER_THRESHOLD_PARAM = "threshold";
const ACCESS_TOKEN_HEADER = "access-token";

module.exports = {
	'calculate-promotion': (event, context) => {
		return new Promise((resolve, reject) => {
			console.log("Starting calculate-promotion lambda function")

			const request = event.extensions.request;
			const threshold = request.query[ORDER_THRESHOLD_PARAM];

			if (validateRequest(request)) {
				const accessToken = request.headers[ACCESS_TOKEN_HEADER];
				const userId = request.query[USER_ID_PARAM];
				const url = `${process.env.EC_SVC_URL}/rest/v2/electronics/users/${userId}/orders`;
				const options = {
					url: url,
					headers: {
						[ACCESS_TOKEN_HEADER]: `Bearer ${accessToken}`
					},
					json: true
				};

				calculatePromotion(options, threshold, resolve, reject);
			} else {
				reject({
					stack: "user-id and threshold params need to be specified"
				})
			}
		})
	}
};

function validateRequest(request) {
	return request.query[USER_ID_PARAM] != undefined && request.query[ORDER_THRESHOLD_PARAM] != undefined;
}

function calculatePromotion(options, threshold, resolve, reject) {
	request.get(options, (error, response, body) => {
		if (!error) {
			if (response.statusCode == 200) {
				console.log('Getting orders from EC succeeded.');
				const s = calculateTotalPrice(body.orders);

				resolve({promotion: s > threshold});
			}
			else {
				reject({
					stack: `Getting orders returned unexpected status: ${response.statusCode}.`
				})
			}
		} else {
			reject({
				stack: "Failed to get orders."
			})
		}
	})
}

function calculateTotalPrice(orders) {
	return orders.reduce((sum, order) => sum + order.total.value, 0)
}
