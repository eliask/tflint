package rules

import (
	"fmt"
	"log"

	"github.com/terraform-linters/tflint/plugin"
	"github.com/terraform-linters/tflint/rules/awsrules"
	"github.com/terraform-linters/tflint/rules/terraformrules"
	"github.com/terraform-linters/tflint/tflint"
)

// Rule is an implementation that receives a Runner and inspects for resources and modules.
type Rule interface {
	Name() string
	Enabled() bool
	Check(runner *tflint.Runner) error
}

// DefaultRules is rules by default
var DefaultRules = append(manualDefaultRules, modelRules...)
var deepCheckRules = append(manualDeepCheckRules, apiRules...)

var manualDefaultRules = []Rule{
	awsrules.NewAwsDBInstanceDefaultParameterGroupRule(),
	awsrules.NewAwsDBInstanceInvalidTypeRule(),
	awsrules.NewAwsDBInstancePreviousTypeRule(),
	awsrules.NewAwsElastiCacheClusterDefaultParameterGroupRule(),
	awsrules.NewAwsElastiCacheClusterInvalidTypeRule(),
	awsrules.NewAwsElastiCacheClusterPreviousTypeRule(),
	awsrules.NewAwsInstancePreviousTypeRule(),
	awsrules.NewAwsRouteNotSpecifiedTargetRule(),
	awsrules.NewAwsRouteSpecifiedMultipleTargetsRule(),
	awsrules.NewAwsS3BucketInvalidACLRule(),
	awsrules.NewAwsS3BucketInvalidRegionRule(),
	awsrules.NewAwsSpotFleetRequestInvalidExcessCapacityTerminationPolicyRule(),
	terraformrules.NewTerraformDashInResourceNameRule(),
	terraformrules.NewTerraformDocumentedOutputsRule(),
	terraformrules.NewTerraformDocumentedVariablesRule(),
	terraformrules.NewTerraformModulePinnedSourceRule(),
	terraformrules.NewTerraformModuleSemverSourceRule(),
}

var manualDeepCheckRules = []Rule{
	awsrules.NewAwsInstanceInvalidAMIRule(),
	awsrules.NewAwsLaunchConfigurationInvalidImageIDRule(),
}

// CheckRuleNames returns map of rules indexed by name
func CheckRuleNames(ruleNames []string, c *tflint.Config) error {
	log.Print("[INFO] Checking rules")

	rulesMap := map[string]Rule{}
	for _, rule := range append(DefaultRules, deepCheckRules...) {
		rulesMap[rule.Name()] = rule
	}

	pluginRules, err := plugin.NewRules(c)
	if err != nil {
		return err
	}
	for _, rule := range pluginRules {
		if _, exists := rulesMap[rule.Name()]; exists {
			return fmt.Errorf("Rule %s is duplicated. Rule names must be unique", rule.Name())
		}
		rulesMap[rule.Name()] = rule
	}

	totalEnabled := 0
	for _, rule := range rulesMap {
		if rule.Enabled() {
			totalEnabled++
		}
	}
	log.Printf("[INFO]   %d (%d) rules total", len(rulesMap), totalEnabled)
	for _, rule := range ruleNames {
		if _, ok := rulesMap[rule]; !ok {
			return fmt.Errorf("Rule not found: %s", rule)
		}
	}
	return nil
}

// NewRules returns rules according to configuration
func NewRules(c *tflint.Config) ([]Rule, error) {
	log.Print("[INFO] Prepare rules")

	ret := []Rule{}
	allRules := []Rule{}

	if c.DeepCheck {
		log.Printf("[DEBUG] Deep check mode is enabled. Add deep check rules")
		allRules = append(DefaultRules, deepCheckRules...)
	} else {
		allRules = DefaultRules
	}

	pluginRules, err := plugin.NewRules(c)
	if err != nil {
		return ret, err
	}
	for _, pluginRule := range pluginRules {
		allRules = append(allRules, pluginRule)
	}

	for _, rule := range allRules {
		enabled := rule.Enabled()
		if r := c.Rules[rule.Name()]; r != nil {
			if r.Enabled {
				log.Printf("[DEBUG] `%s` is enabled", rule.Name())
			} else {
				log.Printf("[DEBUG] `%s` is disabled", rule.Name())
			}
			enabled = r.Enabled
		}

		if enabled {
			ret = append(ret, rule)
		}
	}
	log.Printf("[INFO]   %d rules enabled", len(ret))
	return ret, nil
}
