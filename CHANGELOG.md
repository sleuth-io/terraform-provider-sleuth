## 0.4.1 (July 25, 2023)

FIXES:

- [#103](https://github.com/sleuth-io/terraform-provider-sleuth/pull/103) Translate log

ENHANCEMENTS:

- [#101](https://github.com/sleuth-io/terraform-provider-sleuth/pull/101) Add CUSTOM_INCIDENT as Incident Impact Source
- [#96](https://github.com/sleuth-io/terraform-provider-sleuth/pull/96) Add JIRA provider as Incident Impact Source

## 0.4.0 (July 18, 2023)

ENHANCEMENTS:

- [#87](https://github.com/sleuth-io/terraform-provider-sleuth/pull/87) Add Incident Impact Source resource
- [#94](https://github.com/sleuth-io/terraform-provider-sleuth/pull/94) Add DataDog provider as Incident Impact Source input

## 0.3.11 (July 13, 2023)

FIXES:

- [#83](https://github.com/sleuth-io/terraform-provider-sleuth/pull/83) Environment slug should not show up as updating when re-applying

## 0.3.10 (May 19, 2023)

ENHANCEMENTS:

- Support for Azure in Code Changes resources
- SDK v2 version upgrade

## 0.3.9 (May 11, 2023)

ENHANCEMENTS:

- Documentation improvements, more examples

## 0.3.8 (April 14, 2023)

FIXES:

- create and update of error impact source
- create and update of metric impact source

## 0.3.7 (April 6, 2023)

ENHANCEMENTS:

- Add support for specifying the integration slug for metric impact source
- Add support for specifying whether build branches should match environment branches in build mapping

## 0.3.5 (Sept 27, 2022)

ENHANCEMENTS:

- Add support for specifying the integration slug for metric impact source

## 0.3.4 (Sept 9, 2022)

ENHANCEMENTS:

- Add support for build mapping for code deployments

## 0.3.3 (June 14, 2022)

ENHANCEMENTS:

- Increase HTTP timeout from 10 secs to 20

## 0.3.2 (May 17, 2022)

ENHANCEMENTS:

- Document available deployment tracking choices

## 0.3.1 (April 13, 2022)

ENHANCEMENTS:

- Better handle remote-deleted objects and recreate in that case

## 0.3.0 (March 24, 2022)

ENHANCEMENTS:

- Expand docs about how to use Sleuth resources

NOTES:

- resource/project: The `description` field has been deprecated and will be removed in the next major release

## 0.2.1 (Dec 20, 2021)

ENHANCEMENTS:

- Improve docs, particularly around slugs

## 0.2.0 (Nov 22, 2021)

FEATURES:

- Initial release

## 0.1.0 (Nov 22, 2021)
