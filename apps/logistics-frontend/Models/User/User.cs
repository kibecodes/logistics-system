namespace logistics_frontend.Models.User;
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

public static class UserRoles
{
    public const string Admin = "admin";
    public const string Driver = "driver";
    public const string Customer = "customer";

}

public enum UserStatus
{
    Active,
    Inactive,
    Suspended,
    Pending
}

public class User
{
    [JsonPropertyName("id")]
    public Guid ID { get; set; }

    [JsonPropertyName("fullName")]
    public string FullName { get; set; } = string.Empty;

    [JsonPropertyName("slug")]
    public string Slug { get; set; } = string.Empty;

    [JsonPropertyName("email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("role")]
    public string Role { get; set; } = string.Empty;

    [JsonPropertyName("phone")] 
    public string? Phone { get; set; }

    [JsonPropertyName("token")]
    public string? Token { get; set; }

    [JsonPropertyName("status")]
    [JsonConverter(typeof(JsonStringEnumConverter))]
    public UserStatus? Status { get; set; }

    [JsonPropertyName("last_login")]
    public DateTime? LastLogin { get; set; }
}

public class RegisterModel
{
    [Required(ErrorMessage = "Full name is required")]
    public string? FullName { get; set; }

    [Required(ErrorMessage = "Email is required")]
    [EmailAddress(ErrorMessage = "Email format is invalid")]
    public string? Email { get; set; }

    [Required(ErrorMessage = "Password is required")]
    [MinLength(6, ErrorMessage = "Password must be at least 6 characters")]
    public string? Password { get; set; }

    [Required(ErrorMessage = "Confirm Password is required")]
    [Compare("Password", ErrorMessage = "Passwords do not match")]
    [JsonIgnore]
    public string? ConfirmPassword { get; set; }

    [Required(ErrorMessage = "Role is required")]
    [RegularExpression("admin|driver|customer", ErrorMessage = "Role must be admin, driver, or customer")]
    public string? Role { get; set; }

    [Required(ErrorMessage = "Phone number is required")]
    [Phone(ErrorMessage = "Invalid phone number")]
    public string? Phone { get; set; }
}
