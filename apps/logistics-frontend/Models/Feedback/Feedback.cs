namespace logistics_frontend.Models.Feedback
{
    public class Feedback
    {
        public Guid ID { get; set; }
        public Guid OrderID { get; set; }
        public Guid CustomerID { get; set; }
        public string Comments { get; set; } = string.Empty;
        public DateTime SubmittedAt { get; set; }
    }

    public class CreateFeedbackRequest
    {
        public Guid OrderID { get; set; }
        public Guid CustomerID { get; set; }
        public int Rating { get; set; }
        public string Comment { get; set; } = string.Empty;
    }
}
