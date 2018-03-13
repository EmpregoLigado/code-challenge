package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testitem struct {
	jwt    string
	status int
	reason string
	fail   bool
}

var (
	token   = "aeae42cd8f444313a4f300088713e71c"
	testSet = []testitem{
		//Token gerado com os seguintes dados:
		//payload: {"id": "62a83bce-7caf-455e-a235-57b0ca108b59","exp": 33119884799}
		//header:{"alg": "RS256","typ": "JWT"}
		//Secret aeae42cd8f444313a4f300088713e71c
		//A data de expiração é 31/12/3017
		//Falha devido ao algoritmo errado
		testitem{
			jwt:    "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyYTgzYmNlLTdjYWYtNDU1ZS1hMjM1LTU3YjBjYTEwOGI1OSIsImV4cCI6MzMxMTk4ODQ3OTl9.wVaqqDmMXI040wo2OtCaka76x_KAE5g1o26V9yve-L8",
			status: 401,
			reason: "wrong algorithm",
			fail:   true,
		},

		//Token gerado com os seguintes dados:
		//payload: {"id": "62a83bce-7caf-455e-a235-57b0ca108b59","exp": 33119884799}
		//header:{"alg": "HS256","typ": "JWT"}
		//Secret wrongkey
		//A data de expiração é 31/12/3017
		//Falha devido a chave errada
		testitem{
			jwt:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyYTgzYmNlLTdjYWYtNDU1ZS1hMjM1LTU3YjBjYTEwOGI1OSIsImV4cCI6MzMxMTk4ODQ3OTl9.1SOpzDR7Ucw7U-EJjDTAmigxTh18K3nOHtKXpI-tZD8",
			status: 401,
			reason: "wrong cripto key",
			fail:   true,
		},

		//Token gerado com os seguintes dados:
		//payload: {"id": "62a83bce-7caf-455e-a235-57b0ca108b59","exp": 1499903999}
		//header:{"alg": "HS256","typ": "JWT"}
		//Secret aeae42cd8f444313a4f300088713e71c
		//A data de expiração é 31/12/2015
		//Falha por expiração
		testitem{
			jwt:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYyYTgzYmNlLTdjYWYtNDU1ZS1hMjM1LTU3YjBjYTEwOGI1OSIsImV4cCI6MTQ5OTkwMzk5OX0.zGH8m6473-ak3i2zK1gvZ3LTj156pvNec6NQBiebNO8",
			status: 401,
			reason: "expired key",
			fail:   true,
		},

		//Token gerado com os seguintes dados:
		//payload: {"id": "62a83bce-7caf-455e-a235-57b0ca108b59","exp": 33119884799}
		//header:{"alg": "HS256","typ": "JWT"}
		//Secret aeae42cd8f444313a4f300088713e71c
		//A data de expiração é 31/12/3017
		testitem{
			jwt:    "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjMzMTE5ODg0Nzk5fQ.nHpr7MgyjAIA_I1de-6baw3WU_CvCEuGO54p9Rruqx4",
			status: 200,
			reason: "should have passed, unless it's year 3018",
			fail:   false,
		},
	}
)

func TestAuthHeader(t *testing.T) {
	var h http.Handler
	h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("62a83bce-7caf-455e-a235-57b0ca108b59"))
	})
	jwtsecure := JWTSecure(token)
	ts := httptest.NewServer(jwtsecure(h.ServeHTTP))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatalf(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Fatalf("Should not accept requests without Authorization Header")
	}

	for _, item := range testSet {
		req.Header.Set("Authorization", item.jwt)
		resp, err = client.Do(req)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer resp.Body.Close()
		if resp.StatusCode != item.status {
			t.Fatalf("%s", item.reason)
		} else {
			res, _ := ioutil.ReadAll(resp.Body)
			if item.fail {
				//debug
				//fmt.Println(string(res))
			} else {
				if string(res) != "62a83bce-7caf-455e-a235-57b0ca108b59" {
					t.Fatalf("expected 62a83bce-7caf-455e-a235-57b0ca108b59 got %s", res)
				}
			}
		}
	}
}
