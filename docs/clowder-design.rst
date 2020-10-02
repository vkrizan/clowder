Clowder Design
==============

Problem Statement
-----------------

Engineers have a steep learning curve when they start writing their first
application for cloud.redhat.com, particularly if they're unfamiliar working in
a cloud-native environment before.

Additionally, enforcing standards and consistency across cloud.redhat.com
applications is difficult because of the number of disparate teams across the
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
- Enable engineers to deploy their application in every Kubernetes environment
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

Build a single operator that handles the common use cases that cloud.redhat.com
applications have.  These use cases will be encoded into the API of the
operator, which is of course it CRDs.  There will be two CRDs:

1. ClowdEnvironment

   This CR represents an instance of the entire cloud.redhat.com environment,
   e.g. stage or prod.  It contains configuration for various aspects of the
   environment, implenented by *providers*.

2. ClowdApp

   This CR represents a all the configuration an app needs to be deployed into
   the cloud.redhat.com environment, including:

       - One or more deployment specs
       - Kafka topics
       - Databases
       - Object store buckets (e.g. S3)
       - In-memory DBs (e.g. Redis)
       - Public API endpoint name(s)
       - SLO/SLI thresholds

.. image:: ../images/clowder-flow.svg

Common configuration interface:
.. image:: ../images/clowder-new.svg

Alternatives
------------

- One operator per app with a shared library between them
- One operator with a CRD for each app

.. vim: tw=80
