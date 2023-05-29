namespace FirstService
{
    public class ClientService
    {
        private readonly HttpClient _httpClient;
        private readonly string _url;
        public ClientService(string serviceUrl) { 
            _httpClient = new HttpClient();
            _url = serviceUrl;
        }

        public async Task<string> GetIdFromService()
        {
            try
            {
                using HttpResponseMessage response = await _httpClient.GetAsync(_url);
                response.EnsureSuccessStatusCode();

                string result = await response.Content.ReadAsStringAsync();
                return result;

            }catch (Exception ex) {
                throw new Exception("Error al llamar segundo servicio");
            }
        }
    }
}
