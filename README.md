# imagelist

Imagelist is a simple container image registry with a REST API. Container image
metadata can be added, queried and removed.

The main use case for this is the ability to whitelist images during image build
process and then query them at deployment time. For example, we use it in
Kubernetes with external admission controller - if an image that is about to be
deployed is not found in the imagelist service, the deployment is rejected.

## Getting Started

The service depends on postgres for persistence. More details about various
configuration flags/options can be found by runnning `imagelist --help`.

Auth for `PUT` and `DELETE` requests can be enabled by setting a `--tokens-file
/path/to/tokens.yaml` flag. Imagelist service will watch for any changes to the
tokens file, so the service restart is not necessary. The format of the file is
as follows:

```yaml
# tokens.yaml
tokens:
  - name: foo
    token: 'LF8zuwb2dJMjZLphiSUQHTeA'
    is_disabled: false
```

## API

### Get images

**HTTP Request**

`GET /images`

**Query string parameters**

| Parameter | Description |
| --------- | ----------- |
| `name`    | Get images filtered by image name. |
| `sort`    | Sort results. One of `created_at`, `updated_at`, `name`, `id`. Default: `updated_at`. |

**Response**

`HTTP/1.1 200 OK`

```json
[
  {
    "created_at": "2017-10-27T16:51:06.698062Z",
    "updated_at": "2017-10-27T16:51:06.698062Z",
    "id": "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
    "name": "quay.io/foo/nginx",
    "tags": [
      "latest",
      "v1.0"
    ]
  },
  {
    "created_at": "2017-10-20T15:57:05.948595Z",
    "updated_at": "2017-10-25T10:04:56.426671Z",
    "id": "89562e2c3c4f29abd7e7a3586908f147ee0abeabbf9c07582daf632872b8ac0a",
    "name": "quay.io/bar/proxy",
    "tags": [
      "v1.2.3"
    ]
  }
]
```

### Get an image

**HTTP Request**

`GET /images/:id`

**URL parameters**

| Parameter | Description |
| --------- | ----------- |
| `id`      | Specific image id. |

**Response**

`HTTP/1.1 200 OK`

```json
{
  "created_at": "2017-10-27T16:51:06.698062Z",
  "updated_at": "2017-10-27T16:51:06.698062Z",
  "id": "b5bb9d8014a0f9b1d61e21e796d78dccdf1352f23cd32812f4850b878ae4944c",
  "name": "quay.io/foo/nginx",
  "tags": [
    "latest",
    "v1.0"
  ]
}
```

### Delete an image

**HTTP Request**

`DELETE /images/:id`

Optional: `Authorization: Bearer <token>`

**URL parameters**

| Parameter | Description |
| --------- | ----------- |
| `id`      | Specific image id. |

**Response**

`HTTP/1.1 204 No Content`


### Create or update an image

The same endpoint is used for creating and updating image entries. However,
image tags can be updated only.

**HTTP Request**

`PUT /images`

Optional: `Authorization: Bearer <token>`

**Request data**

```json
{
  "id": "f569d31fa5e46509fe0235f5bc1344fc96c8c142455056c3e02c43dea70f05eb",
  "name": "myorg/apache2",
  "tags": [
    "v0.2.5",
    "latest"
  ]
}
```

**Response**

`HTTP/1.1 200 OK`

```json
{
  "created_at": "2017-10-27T17:07:33.861542Z",
  "updated_at": "2017-10-27T17:07:44.067523Z",
  "id": "f569d31fa5e46509fe0235f5bc1344fc96c8c142455056c3e02c43dea70f05eb",
  "name": "myorg/apache2",
  "tags": [
    "v0.2.5",
    "latest"
  ]
}
```

### Get health

Checks database connection.

**HTTP Request**

`GET /health/ping`

**Response**

`HTTP/1.1 200 OK`


## Development

Project dependencies are managed using [dep](https://github.com/golang/dep).

**Build**

```
$ git clone https://github.com/UKHomeOffice/imagelist.git
$ cd imagelist
$ dep ensure
$ go build
```
