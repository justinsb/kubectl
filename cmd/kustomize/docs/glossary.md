# Glossary

[DAM]: #declarative-application-management
[JSON]: https://www.json.org/
[Resource]: #resource
[YAML]: http://www.yaml.org/start.html
[application]: #application
[apply]: #apply
[apt]: https://en.wikipedia.org/wiki/APT_(Debian)
[base]: #base
[bases]: #base
[bespoke]: #bespoke-configuration
[kustomize]: #kustomize
[kustomization]: #kustomization
[off-the-shelf]: #off-the-shelf
[overlay]: #overlay
[overlays]: #overlay
[patch]: #patch
[patches]: #patch
[proposal]: https://github.com/kubernetes/community/pull/1629
[rebase]: https://git-scm.com/docs/git-rebase
[resource]: #resource
[resources]: #resource
[rpm]: https://en.wikipedia.org/wiki/Rpm_(software)
[target]: #target
[workflow]: workflows.md

## application

An _application_ is a group of k8s resources related by
some common purpose, e.g.  a load balancer in front of a
webserver backed by a database.
[Resource] labelling, naming and metadata schemes have
historically served to group resources together for
collective operations like _list_ and _remove_.

This [proposal] describes a new k8s resource called
_application_ to more formally describe this idea and
provide support for application-level operations and
dashboards.

[kustomize] configures k8s resources, and the proposed
application resource is just another resource.


## apply

The verb _apply_ in the context of k8s refers to a
kubectl command and an in-progress [API
endpoint](https://goo.gl/UbCRuf) for mutating a
cluster.

One _applies_ a statement of what one wants to a
cluster in the form of a complete resource list.

The cluster merges this with the previously applied
state and the actual state to arrive at a new desired
state, which the cluster's reconcilation loop attempts
to create.  This is the foundation of level-based state
management in k8s.

## base

A _base_ is a [target] that some [overlay] modifies.

Any target, including an overlay, can be a base to
another target.

A base has no knowledge of the overlays that refer to it.

A base is usable in isolation, i.e. one should
be able to [apply] a base to a cluster directly.

## bespoke configuration

A _bespoke_ configuration is a [kustomization] and some
[resources] created and maintained internally by some
organization for their own purposes.

The [workflow] associated with a _bespoke_ config is
simpler than the workflow associated with an
[off-the-shelf] config, because there's no notion of
periodically capturing someone else's upgrades to the
[off-the-shelf] config.

## declarative application management

_Declarative Application Management_ (DAM) is a [set of
ideas](https://goo.gl/T66ZcD) aiming to ease management
of k8s clusters.

 * Works with any configuration, be it bespoke,
   off-the-shelf, stateless, stateful, etc.
 * Supports common customizations, and creation of
   instance variants (dev vs, staging vs. production).
 * Exposes and teaches native k8s APIs, rather than
   hiding them.
 * No friction integration with version control to
   support reviews and audit trails.
 * Composable with other tools in a unix sense.
 * Eschews crossing the line into templating, domain
   specific languages, etc., frustrating the other
   goals.

## instance

An _instance_ is the outcome, in a cluster, of applying
an [overlay] to a [base].

> E.g., a _staging_ and _production_ overlay both modify some
> common base to create distinct instances.
>
> The _staging_ instance is the set of resources
> exposed to quality assurance testing, or to some
> external users who'd like to see what the next
> version of production will look like.
>
> The _production_ instance is the set of resources
> exposed to production traffic, and thus may employ
> deployments with a large number of replicas and higher
> cpu and memory requests.

## kustomization

A _kustomization_ is a file called `kustomize.yaml` that
describes a configuration consumable by [kustomize].

Here's an [example](kustomize.yaml).

A kustomization contains fields falling into these categories:

 * Immediate customization instructions -
   _nameprefix_, _labelprefix_, etc.
 * Resource _generators_ for configmaps and secrets.
 * References to _external files_ in these categories:
   * [resources] - completely specified k8s API objects,
      e.g. `deployment.yaml`, `configmap.yaml`, etc.
   * [patches] - _partial_ resources that modify full
     resources defined in a [base]
     (only meaningful in an [overlay]).
   * [bases] - path to a directory containing
     a [kustomization] (only meaningful in an [overlay]).
 * (_TBD_) Standard k8s API kind-version fields.

## kustomize

_kustomize_ is a command line tool supporting template-free
customization of declarative configuration targetted to
k8s.

_Targetted to k8s means_ that kustomize may need some
limited understanding of API resources, k8s concepts
like names, labels, namespaces, etc. and the semantics
of resource patching.

kustomize is an implementation of [DAM].


## off-the-shelf configuration

An _off-the-shelf_ configuration is a kustomization and
resources intentionally published somewhere for others
to use.

E.g. one might create a github repository like this:

> ```
> github.com/username/someapp/
>   kustomize.yaml
>   deployment.yaml
>   configmap.yaml
>   README.md
> ```

Someone could then _fork_ this repo (on github) and
_clone_ their fork to their local disk for
customization.

This clone could act as a [base] for the user's
own [overlays] to do further customization.

## overlay

An _overlay_ is a [target] that modifies (and thus
depends on) another target.

The [kustomization] in an overlay refers to (via file path,
URI or other method) to _some other kustomization_, known as
its [base].

An overlay is unusable without its base.

An overlay supports the typical notion of a
_development_, _QA_, _staging_ and _production_
environment instances.

The configuration of these environments is specified in
individual overlays (one per environment) that all
refer to a common base that holds common configuration.
One configures the cluser like this:

> ```
>  kustomize inflate someapp/overlays/staging |\
>      kubectl apply -f -
>
>  kustomize inflate someapp/overlays/production |\
>      kubectl apply -f -
> ```

Usage of the base is implicit (the overlay's kustomization
points to the base).

An overlay may act as a base to another overlay.

## package

The word _package_ has no meaning in kustomize, as
kustomize is not to be confused with a package
management tool in the tradition of, say, [apt] or
[rpm].

## patch

A _patch_ is a partially defined k8s resource with a
name that must match a resource already known per
traversal rules built into [kustomize].

_Patch_ is a field in the kustomization, distinct from
resources, because a patch file looks like a resource
file, but has different semantics.  A patch depends on
(modifies) a resource, whereas a resourse has no
dependencies.  Since any resource file can be used as a
patch, one cannot reliably distinguish a resource from
a patch just by looking at the file's [YAML].

## resource

A _resource_ is a path to a [YAML] or [JSON] file that
completely defines a functional k8s API object.

## sub-target / sub-application / sub-package

A _sub-whatever_ is not a thing. There are only [bases] and [overlays].

## target

The _target_ is the argument to `build`, e.g.:

> ```
>  kustomize build $target
> ```

`$target` must be a path to a directory that
immediately contains a file called
`kustomize.yaml` (i.e. a [kustomization]).

The target contains, or refers to, all the information
needed to create customized resources to send to the
[apply] operation.

A target is a [base] or an [overlay].
