using Microsoft.AspNetCore.Components;

namespace logistics_frontend.Models.Shared
{
    // Base for storing diff field value types together
    public interface IFieldDefinitionBase
    {
        string Label { get; set; }
        string Type { get; set; } // "text", "select", etc
        List<string> Options { get; set; }
        object? GetValue();
        Task SetValueAsync(object? value);
    }

    // Generic interface for type safety when defining a field
    public interface IFieldDefinition<TValue> : IFieldDefinitionBase
    {
        TValue Value { get; set; }
        EventCallback<TValue> ValueChanged { get; set; }
    }

    // Implementation
    public class FieldDefinition<TValue> : IFieldDefinition<TValue>
    {
        public string Label { get; set; } = "";
        public string Type { get; set; } = "text";
        public List<string> Options { get; set; } = new();
        public TValue Value { get; set; } = default!;
        public EventCallback<TValue> ValueChanged { get; set; }
        public object? GetValue() => Value;

        public async Task SetValueAsync(object? value)
        {
            if (value is TValue typedValue)
                await ValueChanged.InvokeAsync(typedValue);
        }
    }
}
