# bpmn-to-cacao

# About

BPMN-to-CACAO is a utility to convert BPMN XML workflows into [OASIS CACAO](https://www.oasis-open.org/committees/tc_home.php?wg_abbrev=cacao) JSON playbooks.
BPMN-to-CACAO supports CACAO specification [v1.1](https://docs.oasis-open.org/cacao/security-playbooks/v1.1/csd01/security-playbooks-v1.1-csd01.html) and has limited support for [v2.0](https://docs.oasis-open.org/cacao/security-playbooks/v2.0/csd01/security-playbooks-v2.0-csd01.html).

# Compilation

BPMN-to-CACAO requires go1.18 at a minimum. To build the `bpmn-to-cacao` executable:
```
make build
```
To run tests:
```
make test
```

# Usage

Using https://github.com/cisagov/shareable-soar-workflows as input:
```
find ./shareable-soar-workflows -name \*.bpmn -exec bpmn-to-cacao --output-dir=out {} \;
```

# Limitations

This utility is intended to create CACAO playbooks as a starting point.
The conversion from BPMN to CACAO is not perfect, as BPMN is a general purpose notation to specify buiness processes, while CACAO is directly applicable to cybersecurity.
Additionally, BPMN uses "gateways" as an abstraction to support non-linear constructs, while CACAO uses if/while/switch branch statements that are more aligned with procedural programming languages, and BPMN gateways do not always map to if/while/switch statements in a consistent way.
The other limitation is that BPMN workflows may have multiple entry points, implemented either as an explicit "start" action, or using event driven logic (intermediate catch event), while CACAO assumes exactly one start step.
