namespace logistics_frontend.Models.Errors
{
    public class ErrorResponse
    {
        public string Error { get; set; } = string.Empty;
        public string Detail { get; set; } = string.Empty;
    }

    public class ServiceResult<T>
    {
        public bool Success { get; set; }
        public T? Data { get; set; }
        public string? ErrorMessage { get; set; }

        public static ServiceResult<T> Ok(T data) => new() { Success = true, Data = data };
        public static ServiceResult<T> Fail(string message) => new() { Success = false, ErrorMessage = message };
    };
}

