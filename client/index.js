async function register(username) {
  /* Get challenge from server */
  const firstResponse = await fetch(
    `http://localhost:8080/register/start?username=${username}`,
    {
      method: "POST",
    }
  );
  const { publicKey } = await firstResponse.json();
  log("Received server challenge + user data", { publicKey });

  /* Decode server message */
  const decodedResponse = {
    publicKey: {
      ...publicKey,
      challenge: base64UrlByteArray(publicKey.challenge),
      user: {
        ...publicKey.user,
        id: base64UrlByteArray(publicKey.user.id),
      },
    },
  };
  log("Decoded server message", decodedResponse);

  /* Solve challenge + authenticate + create credentials  */
  const credential = await navigator.credentials.create(decodedResponse);
  log("Created keypair using authenticator", credential);

  /* Encode challenge solution */
  const serializedCredential = {
    id: credential.id,
    rawId: arrayBufferToBase64Url(credential.rawId),
    type: "public-key",
    response: {
      clientDataJSON: arrayBufferToBase64Url(
        credential.response.clientDataJSON
      ),
      attestationObject: arrayBufferToBase64Url(
        credential.response.attestationObject
      ),
    },
  };
  log("Encoded response", serializedCredential);

  const secondResponse = await fetch(
    `http://localhost:8080/register/finish?username=${username}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(serializedCredential),
    }
  );
  console.log(await secondResponse.text());
}

async function login(email) {
  /* Get challenge from server */
  const firstResponse = await fetch(
    `http://localhost:8080/login/start?username=${email}`,
    {
      method: "POST",
    }
  );
  const { publicKey } = await firstResponse.json();
  log("Received server challenge + user data", publicKey);

  /* Decode challenge from server */
  const deserializedPublickey = {
    ...publicKey,
    challenge: base64UrlByteArray(publicKey.challenge),
    allowCredentials: publicKey.allowCredentials.map((credential) => ({
      ...credential,
      id: base64UrlByteArray(credential.id),
    })),
  };

  /* Solve challenge + authenticate + get credentials */
  const credential = await navigator.credentials.get({
    publicKey: deserializedPublickey,
  });
  log("Get user credential", credential);

  /* Encode solution */
  const encodedCredential = {
    id: credential.id,
    rawId: arrayBufferToBase64Url(credential.rawId),
    type: "public-key",
    response: {
      ...credential.response,
      authenticatorData: arrayBufferToBase64Url(
        credential.response.authenticatorData
      ),
      clientDataJSON: arrayBufferToBase64Url(
        credential.response.clientDataJSON
      ),
      signature: arrayBufferToBase64Url(credential.response.signature),
      userHandle: arrayBufferToBase64Url(credential.response.userHandle),
    },
  };
  log("Encoded response", encodedCredential);

  const secondResponse = await fetch(
    `http://localhost:8080/login/finish?username=${email}`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(encodedCredential),
    }
  );
  console.log(await secondResponse.text());
}

// Utility function to convert ArrayBuffer to URL-safe Base64 string
function arrayBufferToBase64Url(buffer) {
  return btoa(String.fromCharCode(...new Uint8Array(buffer)))
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");
}

// Utility function to convert URL-safe Base64 string to ByteArray
function base64UrlByteArray(base64Url) {
  const base64 = base64Url
    .replace(/-/g, "+")
    .replace(/_/g, "/")
    .padEnd(base64Url.length + ((4 - (base64Url.length % 4)) % 4), "=");
  return Uint8Array.from(atob(base64), (c) => c.charCodeAt(0));
}

// Utility logging
function log(title, detail) {
  console.groupCollapsed(title);
  console.log(detail);
  console.groupEnd();
}
