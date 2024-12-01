package validation

import (
    "errors"
    "k8s.io/apimachinery/pkg/util/validation/field"
    utilerrors "k8s.io/apimachinery/pkg/util/errors"
    egv1a1 "github.com/envoyproxy/gateway/api/v1alpha1"
)

// ValidateRateLimitRule validates the provided RateLimitRule.
func ValidateRateLimitRule(rule *egv1a1.RateLimitRule) error {
    var errs []error
    if rule == nil {
        return errors.New("rate limit rule is nil")
    }
    if err := validateClientSelectors(rule.ClientSelectors); err != nil {
        errs = append(errs, err)
    }
    if err := validateRateLimitValue(&rule.Limit); err != nil {
        errs = append(errs, err)
    }
    return utilerrors.NewAggregate(errs)
}

// validateClientSelectors validates the ClientSelectors field.
func validateClientSelectors(selectors []egv1a1.RateLimitSelectCondition) error {
    var errs []error
    if len(selectors) > 8 {
        errs = append(errs, field.Invalid(field.NewPath("clientSelectors"), selectors, "must not have more than 8 items"))
    }
    for i, selector := range selectors {
        if err := validateRateLimitSelectCondition(&selector, field.NewPath("clientSelectors").Index(i)); err != nil {
            errs = append(errs, err)
        }
    }
    return utilerrors.NewAggregate(errs)
}

// validateRateLimitSelectCondition validates the RateLimitSelectCondition field.
func validateRateLimitSelectCondition(condition *egv1a1.RateLimitSelectCondition, fldPath *field.Path) error {
    var errs []error
    if len(condition.Headers) == 0 && condition.SourceCIDR == nil {
        errs = append(errs, field.Invalid(fldPath, condition, "at least one of headers or sourceCIDR must be specified"))
    }
    if len(condition.Headers) > 16 {
        errs = append(errs, field.Invalid(fldPath.Child("headers"), condition.Headers, "must not have more than 16 items"))
    }
    for i, header := range condition.Headers {
        if err := validateHeaderMatch(&header, fldPath.Child("headers").Index(i)); err != nil {
            errs = append(errs, err)
        }
    }
    if condition.SourceCIDR != nil {
        if err := validateSourceMatch(condition.SourceCIDR, fldPath.Child("sourceCIDR")); err != nil {
            errs = append(errs, err)
        }
    }
    return utilerrors.NewAggregate(errs)
}

// validateHeaderMatch validates the HeaderMatch field.
func validateHeaderMatch(header *egv1a1.HeaderMatch, fldPath *field.Path) error {
    var errs []error
    if header.Name == "" {
        errs = append(errs, field.Invalid(fldPath.Child("name"), header.Name, "name is required"))
    }
    if len(header.Name) > 256 {
        errs = append(errs, field.Invalid(fldPath.Child("name"), header.Name, "name must not exceed 256 characters"))
    }
    if header.Value != nil && len(*header.Value) > 1024 {
        errs = append(errs, field.Invalid(fldPath.Child("value"), *header.Value, "value must not exceed 1024 characters"))
    }
    return utilerrors.NewAggregate(errs)
}

// validateSourceMatch validates the SourceMatch field.
func validateSourceMatch(source *egv1a1.SourceMatch, fldPath *field.Path) error {
    var errs []error
    if source.Value == "" {
        errs = append(errs, field.Invalid(fldPath.Child("value"), source.Value, "value is required"))
    }
    if len(source.Value) > 256 {
        errs = append(errs, field.Invalid(fldPath.Child("value"), source.Value, "value must not exceed 256 characters"))
    }
    return utilerrors.NewAggregate(errs)
}

// validateRateLimitValue validates the RateLimitValue field.
func validateRateLimitValue(limit *egv1a1.RateLimitValue) error {
    var errs []error
    if limit.Requests == 0 {
        errs = append(errs, field.Invalid(field.NewPath("limit").Child("requests"), limit.Requests, "requests must be greater than 0"))
    }
    return utilerrors.NewAggregate(errs)
}