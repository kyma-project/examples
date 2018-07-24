import jwt from "jsonwebtoken";

import getPem from 'rsa-pem-from-mod-exp';

class AuthService {

    logInByIdToken = async (tokenParams) => {

        if (tokenParams.error) {
            let errMsg = `OIDC provider returned the following error: ${tokenParams.error}.`;
            if (tokenParams.error_description) {
                errMsg += ` Details: ${tokenParams.error_description}`;
            }
            throw Error(errMsg);
        }


        if (!tokenParams.id_token) {
            throw Error("ID token is required!");
        }

        const decodedIdToken = jwt.decode(tokenParams.id_token, {complete: true});

        console.log(`decodedIdToken: ${JSON.stringify(decodedIdToken)}`);

        function userFrom(verifiedIdToken) {
            return {
                id: verifiedIdToken.sub,
                email: verifiedIdToken.email,
                accessToken: tokenParams.access_token,
                idToken: tokenParams.id_token
            };
        }

        const keyID = decodedIdToken.header ? decodedIdToken.header.kid : undefined;
        const issuer = process.env.REACT_APP_OAUTH2_ISSUER;
        const jwksUri = process.env.REACT_APP_OAUTH2_JWKS_URI;

        const jwks = await this.fetchJsonWebKeys(jwksUri);
        const jsonWebKey = this.chooseJsonWebKey(jwks, keyID);
        const key  = this.toKeyValue(jsonWebKey);

        const verificationOptions = {
            audience: 'kyma',
            issuer: issuer,
            ignoreExpiration: false,
            clockTolerance: 5,
            algorithms: [jsonWebKey.alg]
        };

        const verifyToken = new Promise((resolve, reject) => {
            jwt.verify(tokenParams.id_token, key, verificationOptions, function (err, verifiedIdToken) {
                if (err) {
                    return reject(err);
                }
                resolve(verifiedIdToken);
            });
        });
        const verifiedIdToken = await verifyToken;

        return userFrom(verifiedIdToken);
    };

    fetchJsonWebKeys(jwksUri) {
        const requestOptions = {mode: "cors"};
        return fetch(jwksUri, requestOptions).then((response) => response.json())
    }

    chooseJsonWebKey(jwks, keyID) {

        if (jwks && jwks.keys && Array.isArray(jwks.keys)) {

            if (keyID) {

                const found = jwks.keys.filter((key) => {
                    return key.kid === keyID;
                });

                if (found.length > 0) {
                    return found[0];
                }
            }

            return jwks.keys[0];
        }
        return null;
    }

    toKeyValue(jsonWebKey) {

        if (jsonWebKey.value) {
            return jsonWebKey.value;
        }
        else if (jsonWebKey.n && jsonWebKey.e) {
            return getPem(jsonWebKey.n, jsonWebKey.e);
        }
        else {
            return null;
        }
    }
}

export default AuthService;