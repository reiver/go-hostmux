package hostmux


import (
	"testing"

	"net/http"
	"net/http/httptest"
	"strings"
)


func TestNewEquals(t *testing.T) {

	handler := NewEquals()

	if nil == handler {
		t.Errorf("Did not expect NewEquals() func to return nil, but did.")
	}

}


func TestHost(t *testing.T) {

	tests := []struct{
		Hosts         []string
		ExpectedHosts []string
	}{
		{
			Hosts:         []string{},
			ExpectedHosts: []string{},
		},



		{
			Hosts:         []string{"one.com"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"one.com", "two.net"},
			ExpectedHosts: []string{"one.com", "two.net"},
		},
		{
			Hosts:         []string{"one.com", "two.net", "three.org"},
			ExpectedHosts: []string{"one.com", "two.net", "three.org"},
		},



		{
			Hosts:         []string{"one.com", "one.com"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"one.com", "one.com", "one.com"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"one.com", "one.com", "one.com", "one.com"},
			ExpectedHosts: []string{"one.com"},
		},



		{
			Hosts:         []string{"ONE.COM"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"ONE.com"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"one.COM"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"One.Com"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"OnE.cOm"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"oNe.CoM"},
			ExpectedHosts: []string{"one.com"},
		},



		{
			Hosts:         []string{"one.com", "ONE.COM", "ONE.com", "one.COM", "One.Com", "OnE.cOm", "oNe.CoM"},
			ExpectedHosts: []string{"one.com"},
		},
		{
			Hosts:         []string{"one.com", "ONE.COM", "ONE.com", "one.COM", "One.Com", "OnE.cOm", "oNe.CoM",
			                        "two.net", "TWO.NET", "TWO.net", "two.NET", "Two.Net", "TwO.nEt", "tWo.NeT",
			                       },
			ExpectedHosts: []string{"one.com", "two.net"},
		},



		{
			Hosts:         []string{"one.com", "one.com."},
			ExpectedHosts: []string{"one.com", "one.com."},
		},
		{
			Hosts:         []string{"one.com", "one.com.", "two.com", "two.com."},
			ExpectedHosts: []string{"one.com", "one.com.", "two.com", "two.com."},
		},
	}


L:	for testNumber, test := range tests {
		hmux := NewEquals()

		for _, host := range test.Hosts {
			hmux.Host(nil, host)
		}

		if expected, actual := len(test.ExpectedHosts), len(hmux.(*internalEqualsHandler).hostToHandler); expected != actual {
			t.Errorf("For test #%d, expected %d registered host handlers, but got %d.", testNumber, expected, actual)
			continue L
		}

		hadError := false
		for _, expectedHost := range test.ExpectedHosts {
			if _, ok := hmux.(*internalEqualsHandler).hostToHandler[expectedHost]; !ok {
				t.Errorf("For test #%d, expected host %q to be registered, but wasn't.", testNumber, expectedHost)
				hadError = true
			}
		}
		if hadError {
			continue L
		}

		
	}
}


func TestElse(t *testing.T) {

	tests := []struct{
		Hosts []string
		ElseHosts []string
	}{
		{
			Hosts:     []string{},
			ElseHosts: []string{"one.com", "two.net", "three.org"},
		},
		{
			Hosts:     []string{"one.com"},
			ElseHosts: []string{"two.net", "three.org"},
		},
		{
			Hosts:     []string{"one.com", "two.net"},
			ElseHosts: []string{"three.org"},
		},
	}


L:	for testNumber, test := range tests {

		hmux := NewEquals()



		called := map[string]bool{}

		for _, host := range test.Hosts {
			fn := func(w http.ResponseWriter, r *http.Request){
				called[strings.ToLower(r.Host)] = true
			}
			handler := http.HandlerFunc(fn)

			hmux.Host(handler, host)
		}



		elseCalled := map[string]bool{}

		elseFn := func(w http.ResponseWriter, r *http.Request){
			elseCalled[strings.ToLower(r.Host)] = true
		}
		elseHandler := http.HandlerFunc(elseFn)
		hmux.Else(elseHandler)



		hadError := false
		for elseNumber, elseHost := range test.ElseHosts {
			r, err := http.NewRequest("DOES_NOT_MATTER", "/does/not/matter", strings.NewReader("Does not matter."))
			if nil != err {
				t.Errorf("For test #%d and host %q, did not expect an error when creating request. but got one: %v", testNumber, elseHost, err)
				hadError = true
			}
			r.Host = elseHost


			w := httptest.NewRecorder()


			if elseCalled[strings.ToLower(elseHost)] {
				t.Errorf("For test #%d and else host #%d, before triggering expected else host %q not to be already triggered, but was.", testNumber, elseNumber, elseHost)
			}
			for _, host := range test.Hosts {
				if called[strings.ToLower(host)] {
					t.Errorf("For test #%d, before calling ServeHTTP() on hostmux with host %q, a handler not expected to be called was already called.", testNumber, host)
					hadError = true
				}
			}
			hmux.ServeHTTP(w, r)
			for _, host := range test.Hosts {
				if called[strings.ToLower(host)] {
					t.Errorf("For test #%d, after calling ServeHTTP() on hostmux with host %q, a handler not expected to be called was already called.", testNumber, host)
					hadError = true
				}
			}
			if !elseCalled[strings.ToLower(elseHost)] {
				t.Errorf("For test #%d and else host #%d, expected else host %q to trigger else handler, but didn't.", testNumber, elseNumber, elseHost)
			}

		}
		if hadError {
			continue L
		}
	}

}


func TestServeHTTP(t *testing.T) {

	tests := []struct{
		Hosts []string
	}{
		{
			Hosts: []string{"one.com"},
		},
		{
			Hosts: []string{"one.com", "two.net"},
		},
		{
			Hosts: []string{"one.com", "two.net", "three.org"},
		},



		{
			Hosts: []string{"ONE.COM"},
		},
		{
			Hosts: []string{"ONE.COM", "TWO.NET"},
		},
		{
			Hosts: []string{"ONE.COM", "TWO.NET", "THREE.ORG"},
		},
	}


L:	for testNumber, test := range tests {

		hmux := NewEquals()

		called := map[string]bool{}

		for _, host := range test.Hosts {
			fn := func(w http.ResponseWriter, r *http.Request){
				called[strings.ToLower(r.Host)] = true
			}
			handler := http.HandlerFunc(fn)

			hmux.Host(handler, host)
		}


		hadError := false
		for _, host := range test.Hosts {
			r, err := http.NewRequest("DOES_NOT_MATTER", "/does/not/matter", strings.NewReader("Does not matter."))
			if nil != err {
				t.Errorf("For test #%d and host %q, did not expect an error when creating request. but got one: %v", testNumber, host, err)
				hadError = true
			}
			r.Host = host


			w := httptest.NewRecorder()


			if called[strings.ToLower(host)]  {
				t.Errorf("For test #%d, before calling ServeHTTP() on hostmux with host %q, the handler are expecting to be called was already called.", testNumber, host)
			}
			hmux.ServeHTTP(w, r)
			if !called[strings.ToLower(host)]  {
				t.Errorf("For test #%d, when called ServeHTTP() on hostmux with host %q, the handler expected to get called did not get called.", testNumber, host)
			}

		}
		if hadError {
			continue L
		}

	}
}
