Clowder Design
==============

Problem Statement
-----------------

Engineers have a steep learning curve when they start writing their first
application for cloud.redhat.com, particularly if they're unfamiliar working in
a cloud-native environment before.

Additionally, enforcing standards and consistency across cloud.redhat.com
applications is difficult because of the number of disparate teams arcross the
enterprise that contribute to the ecosystem.

Similar to application consistency, it is difficult for engineers to be able to
handle the inconsistencies between environments, e.g local development versus
production.  While these differences can never be fully eliminated, ideally
developers would not have to devise their own solutions to handle these
differences.

Mission
-------

Clowder aims to abstract the underlying Kubernetes environment and configuration
to simplify the creation and maintenance of applications on cloud.redhat.com.

Goals
-----

- Abstract Kubernetes environment and configuration details from applications
- Enable engineers to deploy their appliation in every Kubernetes environmet
  (i.e. dev, stage, prod) with minimal changes to their custom resources.
- Increase operational consistency between cloud.redhat.com applications,
  including rollout parameters and pod affinity.
- Handle metrics configuration and standard SLI/SLO metrics and alerts
- Some form of progressive deployment

Non-goals
---------

- Support applications outside cloud.redhat.com

Proposal
--------

.. image:: ../images/clowder-flow.svg

Alternatives
------------

- One operator per app with a shared library between them
- One operator with a CRD for each app

.. vim: tw=80
