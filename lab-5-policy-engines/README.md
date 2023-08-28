# Policy Engine Example - Open Policy Agent (OPA)

## OPA and Rego

OPA is a general purpose policy engine that can make policy decisions on structured data using the Rego language. We will explore some Rego features within this lab, but the [Policy Language](https://www.openpolicyagent.org/docs/latest/policy-language/) and [Policy Reference](https://www.openpolicyagent.org/docs/latest/policy-reference/) pages on the OPA site can be consulted for more information.

## Example Scenario

We will use an example in this lab where policy decisions are required to ensure that only authorized people within a company can view employee records. The following statements describe the business logic:

- an administrator user can view any employee's record
- an employee's direct line manager can view their record

We will abstract away details of where the policy decision will be enforced (this could be via a call to OPA from within the application itself, or perhaps via a network proxy filtering traffic before it reaches the application). All we need to know for this lab is that the following information is present in an HTTP request to the application, for example in a header or the body of the request:

- a `bearer` field containing a JWT including claims about the authenticated calling user, including their name and whether they are an application administrator
- an `employee_id` representing the employee whose record is being accessed
- an `action` to be performed on the record, e.g. read, update, delete. In this simple lab we will not make a distinction between different actions.

This data will be made available to OPA as `input`, as we shall see later.

### Obtaining an Example JWT

As we are abstracting away details regarding user authentication, we need to be able to create example JWTs for this lab. We will introduce a few features of OPA by performing this task in Rego, using the inbuilt `io.jwt.encode_sign` function (for more information on inbuilt functions, see [Policy Reference](https://www.openpolicyagent.org/docs/latest/policy-reference/)). Look through the Rego in the following file: [create_jwt.rego](create_jwt.rego)

Note that the file starts with a hierarchical package name, so that Rego rules and policies can be organized based on their functionalty. We then see the use of the [ceil](https://www.openpolicyagent.org/docs/latest/policy-reference/#builtin-numbers-ceil) and [time.now_ns](https://www.openpolicyagent.org/docs/latest/policy-reference/#builtin-time-timenow_ns) built in functions, which are used to define an expiry time for our JWTs, one day (don't do this outside of a test environment - JWTs should be short lived!) from the time at which the JWT is created. Note the use of the assignment operator `:=` in the definition of `expiry_time`. Rego supports three kinds of equality:

- assignment (`:=`) - assigned variables are locally scoped within rules
- comparison (`==`) - does recursive, semantic equality checks between values within a rule
- unification (`=`) - combines assignment and comparison. Rego assigns as many variables as it needs to in order to make the comparison true, e.g. `[1,x]=[y,2]` assigns `x` to 2 and `y` to 1.

Once we have an `expiry_time` for our tokens, we can create JWTs for Alice and Bob using the `io.jwt.encode_sign` function. Note that we have an RSA key pair in this file - this has been done for example purposes and is not safe for production usage as this data contains the RSA private key!

Try running OPA in a Docker container to obtain Alice's JWT token using the [`opa eval`](https://www.openpolicyagent.org/docs/latest/#2-try-opa-eval) sub-command. Note that you will require [jq](https://jqlang.github.io/jq/download/) to be installed to run this command.

```bash
docker run -v .:/example openpolicyagent/opa eval -d /example/create_jwt.rego 'data.example.jwt.alice_token' | jq '.result[0].expressions[0].value'
```

Note that the `create_jwt.rego` file has been loaded into OPA using the `-d` / `--data` flag, and we are asking OPA for the value of `data.example.jwt.alice_token`. This shows how the heirarchical package and rules sit under OPA's `data` [Document](https://www.openpolicyagent.org/docs/latest/philosophy/#the-opa-document-model).

Copy the JWT output from the above command (without the surrounding `"`s), and try pasting it into the 'Encoded' box on [jwt.io](https://jwt.io/) to see the decoded claims that match what we defined in [create_jwt.rego](create_jwt.rego). As we have an `is_admin` claim with the JWT, we can implement the first business rule by decoding the JWT, checking it has been signed by a trusted authority, and checking whether our user is an admin. However, to find out whether the calling user is someone's line manager, we need external data.

### External Data

We have considered how external data can be provided to OPA via JWT tokens and other input. However, when external data changes infrequently and can reasonably be stored in memory all at once (such as employee line management data in our example), it can be replicated in bulk via OPA’s [bundle](https://www.openpolicyagent.org/docs/latest/external-data/#option-3-bundle-api) feature. Policies and external data can both be added to a bundle (which is a `.tar.gz` file), which OPA can then consume via a bundle server. Look at the structure of the [bundle](bundle/) directory to see where our external data will sit within OPA's `data` document. As we have a [`data.json`](bundle/user_data/data.json) file within a `user_data` directory within our top-level bundle directory, the information regarding our example users will be found at `data.user_data.users`.

### Writing Policies

We are now in a position to write the policy that will implement our business rules. Take a look at the [policy file](bundle/example/authz/policy.rego). The result of our policy decision will be captured in the value of `allow`. Given what we know about the hierarchical nature of packages and the OPA `data` document, this information will be available at `data.example.authz.allow`.

The input that we will provide to OPA will be in the form:

```json
{ "input": { "bearer": "<JWT here>", "action": "read", "employee_id": 4 } }
```

We can refer to `input.bearer`, `input.action` and `input.employee_id` within our policy, but adding an `import` statement at the top of the policy, e.g. `import input.bearer`, means that we can simply refer to `bearer` in the Rego rules within.

OPA policies are formed from a collection of rules, where rules can take the form `assignment if { conditions}`, e.g.

```rego
allow {
 token_is_valid
 user_is_admin
}
```

The rule body between `{}` is a collection of assignments and expressions. `allow` will evaluate to true if a logical AND of all the assignments and expressions is `true`. If an assignment is false or undefined, `allow` is also undefined. As such, we need to set `default allow := false` in the policy, so that `allow` can only be `true` if one of the rules evaluates to `true` - otherwise it will be `false`, but never undefined. In this way, multiple `allow` rules represent a logical OR.

Read through and understand the rest of the policy, referring out to the [OPA documentation](https://www.openpolicyagent.org/docs/latest/) if necessary, for example to understand iteration using the [`some`](https://www.openpolicyagent.org/docs/latest/) keyword.

### Running the example

Build the policy and data bundle:

```bash
docker run -v .:/example openpolicyagent/opa build \
    --bundle /example/bundle \
    -o example/bundle.tar.gz
```

Run OPA as a server in a Docker container, similar to when we created an example JWT:

```bash
docker run --rm --network host -v .:/example \
    openpolicyagent/opa run --server \
    --bundle /example/bundle.tar.gz \
    --addr localhost:8181
```

The server is now listening on `http://localhost:8181`

Switch to a second terminal tab, ensuring you are in the same directory, and check the results of the following policy decisions:

- can Alice view Eve's record?
- can Bob view Charlie's record?
- can Bob view Dan's record?

It will be convenient to have example JWTs for Alice and Bob exported as environment variables for this task, so run the following two commands (noting from [create_jwt.rego](create_jwt.rego) that we have set the JWT to have a 1 day expiry to give you plenty of time to complete the lab without rerunning the commands):

```bash
export ALICE_JWT=$(docker run -v .:/example openpolicyagent/opa eval -d /example/create_jwt.rego 'data.example.jwt.alice_token' | jq '.result[0].expressions[0].value')
```

```bash
export BOB_JWT=$(docker run -v .:/example openpolicyagent/opa eval -d /example/create_jwt.rego 'data.example.jwt.bob_token' | jq '.result[0].expressions[0].value')
```

In order to check the policy decisions, we can make POST requests to the OPA server of the following form:

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$ALICE_JWT"', "action": "read", "employee_id": 5}}' localhost:8181/v1/data/example/authz/allow
```

Note we are providing the relevant `input` via the parameters in the body of the POST request, and we are querying the value of `data.example.authz.allow` via OPA's [Data API](https://www.openpolicyagent.org/docs/latest/rest-api/#data-api).

To check the three example scenarios above, you can run the following three commands from your second terminal tab:

Can Alice view Eve's record?

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$ALICE_JWT"', "action": "read", "employee_id": 5}}' \
    localhost:8181/v1/data/example/authz/allow
```

Can Bob view Charlie's record?

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$BOB_JWT"', "action": "read", "employee_id": 3}}' \
    localhost:8181/v1/data/example/authz/allow
```

Can Bob view Dan's record?

```bash
curl -X POST -H "Content-Type: application/json" \
    -d '{"input": {"bearer": '"$BOB_JWT"', "action": "read", "employee_id": 4}}' \
    localhost:8181/v1/data/example/authz/allow
```

### Optional Challenge Step

You should notice that whilst Alice can view Eve's record (as Alice is an admin), and Bob can view Charlie's (as Bob directly manages Charlie), Bob can not view Dan's record, as even though he manages Bob's manager (Charlie), our policy only covers direct line managers being able to view employee records. Try to modify the policy to include a rule that would allow a manager to view the records of a person two layers down in the heirarchy. I.e. Bob should be able to view Dan's record, as Bob manages Charlie, and Charlie manages Dan.

### Teardown

Type `Ctrl + C` in the terminal tab running the OPA server, and run:

```bash
rm bundle.tar.gz
```
