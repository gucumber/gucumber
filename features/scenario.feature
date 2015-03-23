Feature: Scenarios
  This package should support the use of scenarios.

  Scenario: Basic usage
    Given I have an initial step
    And I have a second step
    When I run the "cucumber.go" command
    Then this scenario should execute 1 time and pass
