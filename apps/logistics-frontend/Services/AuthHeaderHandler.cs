using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading;
using System.Threading.Tasks;

namespace logistics_frontend.Services.AuthHeaderHandler;

public class AuthHeaderHandler : DelegatingHandler
{
    private readonly UserSessionService _session;
    public AuthHeaderHandler(UserSessionService session)
    {
        _session = session;
    }

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken)
    {
        var token = await _session.GetTokenAsync();

        if (!string.IsNullOrWhiteSpace(token))
        {
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", token);
        }

        return await base.SendAsync(request, cancellationToken);
    }
}

