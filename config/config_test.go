package config

import (
	"fmt"
	"github.com/bivas/rivi/config/action"
	"github.com/bivas/rivi/config/client"
	"github.com/stretchr/testify/assert"
	. "testing"
)

func assertClientConfig(t *T, config client.ClientConfig) {
	assert.Equal(t, config.GetOAuthToken(), "github-token", "oauth token")
	assert.Equal(t, config.GetSecret(), "github-secret", "secret")
}

func assertRoles(t *T, configuration Configuration) {
	roles := configuration.GetRoles()
	assert.Contains(t, roles, "admins", "roles")
	assert.Contains(t, roles, "reviewers", "roles")
	assert.Contains(t, roles, "testers", "roles")
	assert.NotContains(t, roles, "dummy", "dummy role")

	admins := configuration.GetRoleMembers("admins")
	assert.Contains(t, admins, "user1", "admin.user1")
	assert.Contains(t, admins, "user2", "admin.user2")
	assert.NotContains(t, admins, "user3", "admin.user3")

	reviewers := configuration.GetRoleMembers("reviewers")
	assert.Contains(t, reviewers, "user1", "reviewers.user1")
	assert.NotContains(t, reviewers, "user2", "reviewers.user2")
	assert.Contains(t, reviewers, "user3", "reviewers.user3")

	testers := configuration.GetRoleMembers("testers")
	assert.Contains(t, testers, "user2", "testers.user2")
	assert.NotContains(t, testers, "user1", "reviewers.user3")
	assert.NotContains(t, testers, "user3", "reviewers.user3")
}

func assertRules(t *T, configuration Configuration) {
	rules := configuration.GetRules()
	assert.Len(t, rules, 4, "rules len")
	ruleNames := make([]string, 0)
	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name())
	}
	assert.Contains(t, ruleNames, "rule1", "rule name")
	assert.Contains(t, ruleNames, "rule2", "rule name")
	assert.Contains(t, ruleNames, "rule3", "rule name")
	assert.Contains(t, ruleNames, "rule4", "rule name")
	assert.Equal(t, "rule4", ruleNames[0], "first rule")
	assert.Equal(t, "rule3", ruleNames[1], "second rule")
	assert.Equal(t, "rule2", ruleNames[2], "third rule")
	assert.Equal(t, "rule1", ruleNames[3], "fourth rule")
}

func TestReadConfig(t *T) {
	c, err := NewConfiguration("config_test.yml")
	if err != nil {
		t.Fatalf("Got error during config read. %s", err)
	}
	assertClientConfig(t, c.GetClientConfig())
	assertRoles(t, c)
	assertRules(t, c)
}

type testActionConfig struct {
	key, value string
}

func (t *testActionConfig) Name() string {
	return "test-section"
}

type testBuilder struct {
}

func (*testBuilder) Build(config map[string]interface{}) (action.ActionConfig, error) {
	if len(config) != 1 {
		return nil, fmt.Errorf("Wrong number of values")
	}
	for key, value := range config {
		return &testActionConfig{key, value.(string)}, nil
	}
	return nil, fmt.Errorf("Should not reach here")
}

func TestActionConfigBuilder(t *T) {
	action.RegisterActionConfigBuilder("test-section", &testBuilder{})
	c, err := NewConfiguration("config_test.yml")
	if err != nil {
		t.Fatalf("Got error during config read. %s", err)
	}
	result, err := c.GetActionConfig("test-section")
	assert.Nil(t, err, "should contain section")
	exact, ok := result.(*testActionConfig)
	assert.True(t, ok, "should be of test action config type")
	assert.Equal(t, "key", exact.key, "key")
	assert.Equal(t, "value", exact.value, "value")
}
