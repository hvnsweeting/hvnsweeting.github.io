# Django as a library
(Django for people come from Flask, or pure Python).

I've tried to learn Django so many times, IIRC, from Django 1.4, and now it is 2.0
and I still haven't got it. I just follow the main tutorial on its homepage,
finished it for so many times, but that does not make me understand Django
at all. Something was wrong, and now I think I know it.
They way I "tried" Django is not the way I used Python libraries.
For Django, it is top-down, and for all the rest, I learn by bottom-up
approach.

## What is the differences
By learning bottom-up, I know functions, classes that I have, what attributes
they have, and build things from them. Says, Flask, it is also a web-framework,
     but it is very Python-like. A Flask view returns a Python string, that's it.
In Django, the tutorial throws ton of methods, class, shortcuts to you,
   you might even cannot answer on top of your head what a view returns, and
   constantly looking up into the document/tutorial, over and over again.

Let's learn Django in a different way. Treat it as a library.

### Requests and Responses

When an user accesses the website by using a web browser, the web browser sends
an HTTP request to the website.
Django will receive that HTTP request, represents it by an `HttpRequest` object.
Django router will handle the request by a view which associates with
the requested URL.

The view (often a function, or a class with appropriate method), will be called
with the requests object as first argument, process it and return to the web-browser a HTTP response.
A HTTP response contains some data:
- A HTTP status code indicates the result status - OK, Redirected, Not Found ...
- A text (Python string or byte in Python3) - so called `content`, which is usually HTML code, which
will be rendered by browser and displays beautifully.

Django represents a HTTP response by `HttpResponse` object.

Both HTTP request and response can be seen by "console" or "network" on
modern (2017+) web browsers. Most of them use Ctrl+U (View source) to view the text returned in response.

All HTTP Request/Response classes are in `django.http` sub-module.

Let's create an `HttpResponse` objects.

Written used Django at version

```Python
In [22]: django.VERSION
Out[22]: (1, 11, 2, 'final', 0)
```

```Python
In [19]: classes(django.http)
Out[19]: ['BadHeaderError', 'FileResponse', 'Http404', 'HttpRequest', 'HttpResponse', 'HttpResponseBadRequest', 'HttpResponseForbidden', 'HttpResponseGone', 'HttpResponseNotAllowed', 'HttpResponseNotFound', 'HttpResponseNotModified', 'HttpResponsePermanentRedirect', 'HttpResponseRedirect', 'HttpResponseServerError', 'JsonResponse', 'QueryDict', 'RawPostDataException', 'SimpleCookie', 'StreamingHttpResponse', 'UnreadablePostError']

In [20]: response = django.http.HttpResponse('This is content')

In [21]: vars(response)
Out[21]: {'cookies': <SimpleCookie: >, 'closed': False, '_handler_class': None, '_charset': None, '_closable_objects': [], '_container': [b'This is content'], '_reason_phrase': None, '_headers': {'content-type': ('Content-Type', 'text/html; charset=utf-8')}}
```

`vars` is very useful, it prints all data attributes of given object.

So the response object contains some following data:
- cookies (small string of data, used by browser/web app). Current cookie
is empty.
- `_headers`: contains raw data about HTTP headers, often contains type of response (a file? a normal text? a JSON data? which charset?)
- `_container`: contains raw data (in a Python list), which will be join later
and create response content/text. This is what content user will see.
- some other attributes that use for other purposes.

So notice the attribute names which start with underscore (`_`), that says "please do not change this directly", use other methods. Because in Python, you can change almost every mutable thing.

All "public" (name not start with underscore) attributes (includes methods) of response:

```Python
In [42]: pub(response)
Out[42]: ['charset', 'close', 'closed', 'content', 'cookies', 'delete_cookie', 'flush', 'get', 'getvalue', 'has_header', 'items', 'make_bytes', 'readable', 'reason_phrase', 'seekable', 'serialize', 'serialize_headers', 'set_cookie', 'set_signed_cookie', 'setdefault', 'status_code', 'streaming', 'tell', 'writable', 'write', 'writelines']

```

If you are familiar with Python files, these methods `flush`, `write`, `writelines`, `tell`... are to make
response object behave as a file.

Those interesting methods are:
- `status_code`: HTTP status code (int)
- `reason_phrase`: corresponding reason phrase of the `status_code` (str)

```Python
In [43]: response.status_code, response.reason_phrase
Out[43]: (200, 'OK')
```

- `has_header`: check if a header exist

```Python
In [45]: response._headers
Out[45]: {'content-type': ('Content-Type', 'text/html; charset=utf-8')}

In [46]: response.has_header('content-type')
Out[46]: True

In [47]: response.has_header('hack')
Out[47]: False
```

Manipulate cookie:

- `delete_cookie`
- `set_cookie`
- `set_signed_cookie`

Let update cookie and see:

```Python
In [52]: response.set_cookie('username', 'HVN')

In [53]: response.cookies
Out[53]: <SimpleCookie: username='HVN'>

In [55]: response.set_signed_cookie('password', '123456')

In [57]: response.cookies
Out[57]: <SimpleCookie: password='123456:1e9Rns:I6L7NFsJiOgRVfBqAGSoHmTXE3I' username='HVN'>

In [58]: response.delete_cookie('username'); response.cookies
Out[58]: <SimpleCookie: password='123456:1e9Rns:I6L7NFsJiOgRVfBqAGSoHmTXE3I' username=''>
```
So in Django, a view must return an instance of `HttpResponse` or its sub-class instance.

Some sub-class are pre-defined for using to return simple responses without content:

```
In [59]: classes(django.http.response)
Out[59]: ['BadHeaderError', 'DisallowedRedirect', 'DjangoJSONEncoder', 'FileResponse', 'Header', 'Http404', 'HttpResponse', 'HttpResponseBadRequest', 'HttpResponseBase', 'HttpResponseForbidden', 'HttpResponseGone', 'HttpResponseNotAllowed', 'HttpResponseNotFound', 'HttpResponseNotModified', 'HttpResponsePermanentRedirect', 'HttpResponseRedirect', 'HttpResponseRedirectBase', 'HttpResponseServerError', 'JsonResponse', 'SimpleCookie', 'StreamingHttpResponse', 'map']
```

Namely:
- HttpResponseBadRequest:

```Python
class HttpResponseBadRequest(HttpResponse):
    status_code = 400
```

- HttpResponseServerError:

```
class HttpResponseServerError(HttpResponse):
    status_code = 500
```

- HttpResponseRedirect:

```
class HttpResponseRedirect(HttpResponseRedirectBase):
    status_code = 302
```

They just simply set class data attribute `status_code`.

Some non-trivial classes:
- StreamingHttpResponse: only when you want to streaming many data (big file for example).
- JsonResponse: this is so handy if you want to return a JSON response, the class sets header properly.
```Python
class JsonResponse(HttpResponse):
    """
    An HTTP response class that consumes data to be serialized to JSON.

    :param data: Data to be dumped into json. By default only ``dict`` objects
      are allowed to be passed due to a security flaw before EcmaScript 5. See
      the ``safe`` parameter for more information.
    :param encoder: Should be an json encoder class. Defaults to
      ``django.core.serializers.json.DjangoJSONEncoder``.
    :param safe: Controls if only ``dict`` objects may be serialized. Defaults
      to ``True``.
    :param json_dumps_params: A dictionary of kwargs passed to json.dumps().
    """

    def __init__(self, data, encoder=DjangoJSONEncoder, safe=True,
                 json_dumps_params=None, **kwargs):
        if safe and not isinstance(data, dict):
            raise TypeError(
                'In order to allow non-dict objects to be serialized set the '
                'safe parameter to False.'
            )
        if json_dumps_params is None:
            json_dumps_params = {}
        kwargs.setdefault('content_type', 'application/json')
        data = json.dumps(data, cls=encoder, **json_dumps_params)
        super(JsonResponse, self).__init__(content=data, **kwargs)
```

The code is very simple, it set appropriate `content_type`, then json.dump your given data dict - it is a security problem that it avoids to
dump non-dict object.

```Python
In [73]: jr = django.http.JsonResponse({'course': 'PyMivn'})

In [74]: vars(jr)
Out[74]: {'cookies': <SimpleCookie: >, 'closed': False, '_handler_class': None, '_charset': None, '_closable_objects': [], '_container': [b'{"course": "PyMivn"}'], '_reason_phrase': None, '_headers': {'content-type': ('Content-Type', 'application/json')}}
```

All the details can see in
https://docs.djangoproject.com/en/1.11/ref/request-response/
and the code is surprising easy to read.

### Request
### TBD
