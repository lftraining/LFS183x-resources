package envoy.authz

import input.attributes.request.http.method
import input.attributes.request.http.path
import input.attributes.request.http.headers.authorization as bearer

default allow := false

allow {
	token_is_valid
	user_is_admin
}

allow {
	token_is_valid
	user_is_manager
}

jwks := {
    "keys": [{
        "kty":"RSA",
        "n": "ofgWCuLjybRlzo0tZWJjNiuSfb4p4fAkd_wWJcyQoTbji9k0l8W26mPddxHmfHQp-Vaw-4qPCJrcS2mJPMEzP1Pt0Bm4d4QlL-yRT-SFd2lZS-pCgNMsD1W_YpRPEwOWvG6b32690r2jZ47soMZo9wGzjb_7OMg0LOL-bSf63kpaSHSXndS5z5rexMdbBYUsLA9e-KXBdQOS-UTo7WTBEMa2R2CapHg665xsmtdVMTBQY4uDZlxvb3qCo5ZwKh9kG4LT6_I5IhlJH7aGhyxXFvUK-DWNmoudF8NAco9_h9iaGNj8q2ethFkMLs91kzk2PAcDTW9gb54h4FRWyuXpoQ",
        "e": "AQAB",
		"alg": "RS256",
		"use": "sig",
		"kid": "1"
    }]
}

bearer_jwt := authz {
    split_token := split(bearer, "Bearer ")
	authz := split_token[1]
}

token_is_valid := valid {
	[valid, header, payload] := io.jwt.decode_verify(bearer_jwt, {"cert": json.marshal(jwks), "aud": "employee-records"})
}

token_payload := payload {
    [header, payload, signature] := io.jwt.decode(bearer_jwt)
}

user_is_admin := admin {
	admin := token_payload.is_admin
}

specific_user_request {
	method == "GET"
    glob_url := "/api/employees/*"
    glob.match(glob_url, ["/"], path)
}

user_is_manager {
    specific_user_request
    split_path := split(path, "/")
    e_id := split_path[3]
	username := token_payload.sub
	some i
	reportee := data.user_data.users[username].reportees[i]
	format_int(data.user_data.users[reportee].employee_id, 10) == e_id
}