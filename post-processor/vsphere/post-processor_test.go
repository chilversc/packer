package vsphere

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func getTestConfig() Config {
	return Config{
		Username:   "me",
		Password:   "notpassword",
		Host:       "myhost",
		Datacenter: "mydc",
		Cluster:    "mycluster",
		VMName:     "my vm",
		Datastore:  "my datastore",
		Insecure:   true,
		DiskMode:   "thin",
		VMFolder:   "my folder",
	}
}

func TestArgs(t *testing.T) {
	var p PostProcessor

	p.config = getTestConfig()

	source := "something.vmx"
	ovftool_uri := fmt.Sprintf("vi://%s:%s@%s/%s/host/%s",
		url.QueryEscape(p.config.Username),
		url.QueryEscape(p.config.Password),
		p.config.Host,
		p.config.Datacenter,
		p.config.Cluster)

	if p.config.ResourcePool != "" {
		ovftool_uri += "/Resources/" + p.config.ResourcePool
	}

	args, err := p.BuildArgs(source, ovftool_uri)
	if err != nil {
		t.Errorf("Error: %s", err)
	}

	t.Logf("ovftool %s", strings.Join(args, " "))
}

func TestGenerateURI_Basic(t *testing.T) {
	var p PostProcessor

	p.config = getTestConfig()

	uri := p.generateURI()
	expected_uri := "vi://me:notpassword@myhost/mydc/host/mycluster"
	if uri != expected_uri {
		t.Fatalf("URI did not match. Recieved: %s. Expected: %s", uri, expected_uri)
	}
}

func TestGenerateURI_PasswordEscapes(t *testing.T) {
	type escapeCases struct {
		Input    string
		Expected string
	}

	cases := []escapeCases{
		{`this has spaces`, `this%20has%20spaces`},
		{`exclaimation_!`, `exclaimation_%21`},
		{`hash_#_dollar_$`, `hash_%23_dollar_%24`},
		{`ampersand_&awesome`, `ampersand_%26awesome`},
		{`single_quote_'_and_another_'`, `single_quote_%27_and_another_%27`},
		{`open_paren_(_close_paren_)`, `open_paren_%28_close_paren_%29`},
		{`asterisk_*_plus_+`, `asterisk_%2A_plus_%2B`},
		{`comma_,slash_/`, `comma_%2Cslash_%2F`},
		{`colon_:semicolon_;`, `colon_%3Asemicolon_%3B`},
		{`equal_=question_?`, `equal_%3Dquestion_%3F`},
		{`at_@`, `at_%40`},
		{`open_bracket_[closed_bracket]`, `open_bracket_%5Bclosed_bracket%5D`},
		{`user:password with $paces@host/name.foo`, `user%3Apassword%20with%20%24paces%40host%2Fname.foo`},
	}

	for _, escapeCase := range cases {
		var p PostProcessor

		p.config = getTestConfig()
		p.config.Password = escapeCase.Input

		uri := p.generateURI()
		expected_uri := fmt.Sprintf("vi://me:%s@myhost/mydc/host/mycluster", escapeCase.Expected)

		if uri != expected_uri {
			t.Fatalf("URI did not match. Recieved: %s. Expected: %s", uri, expected_uri)
		}
	}
}

func TestEscaping(t *testing.T) {
	type escapeCases struct {
		Input    string
		Expected string
	}

	cases := []escapeCases{
		{`this has spaces`, `this%20has%20spaces`},
		{`exclaimation_!`, `exclaimation_%21`},
		{`hash_#_dollar_$`, `hash_%23_dollar_%24`},
		{`ampersand_&awesome`, `ampersand_%26awesome`},
		{`single_quote_'_and_another_'`, `single_quote_%27_and_another_%27`},
		{`open_paren_(_close_paren_)`, `open_paren_%28_close_paren_%29`},
		{`asterisk_*_plus_+`, `asterisk_%2A_plus_%2B`},
		{`comma_,slash_/`, `comma_%2Cslash_%2F`},
		{`colon_:semicolon_;`, `colon_%3Asemicolon_%3B`},
		{`equal_=question_?`, `equal_%3Dquestion_%3F`},
		{`at_@`, `at_%40`},
		{`open_bracket_[closed_bracket]`, `open_bracket_%5Bclosed_bracket%5D`},
		{`user:password with $paces@host/name.foo`, `user%3Apassword%20with%20%24paces%40host%2Fname.foo`},
	}
	for _, escapeCase := range cases {
		received := escapeWithSpaces(escapeCase.Input)

		if escapeCase.Expected != received {
			t.Errorf("Error escaping URL; expected %s got %s", escapeCase.Expected, received)
		}
	}

}
