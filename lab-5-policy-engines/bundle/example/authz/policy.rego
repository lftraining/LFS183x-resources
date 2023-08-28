package example.authz

import input.bearer

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

token_is_valid := valid {
	[valid, header, payload] := io.jwt.decode_verify(bearer, {"cert": json.marshal(jwks), "aud": "employee-records"})
}

user_is_admin := admin {
	[header, payload, signature] := io.jwt.decode(bearer)
	admin := payload.is_admin
}

user_is_manager {
	[header, payload, signature] := io.jwt.decode(bearer)
	username := payload.sub
	some i
	reportee := data.user_data.users[username].reportees[i]
	data.user_data.users[reportee].employee_id == input.employee_id
}