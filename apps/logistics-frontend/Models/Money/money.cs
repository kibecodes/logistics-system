using System.Text.Json.Serialization;

public class Money
{
    [JsonPropertyName("amount")]
    public long Amount { get; set; } // Smallest unit (e.g., cents)

    [JsonPropertyName("currency")]
    public string Currency { get; set; } = "KES"; // ISO code
}
