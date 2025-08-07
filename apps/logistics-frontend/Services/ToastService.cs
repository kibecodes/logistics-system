public class ToastService
{
    public event Action<string, ToastLevel>? OnShow;

    public void ShowToast(string message, ToastLevel level = ToastLevel.Info)
    {
        OnShow?.Invoke(message, level);
    }
    public enum ToastLevel
    {
        Info,
        Success,
        Warning,
        Error
    }
}