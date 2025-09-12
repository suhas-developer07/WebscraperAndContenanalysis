import axios from "axios";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
const SEARCH_BASE_URL =
  import.meta.env.VITE_SEARCH_BASE_URL || "http://localhost:17029";

export const api = {
  submitUrls: async (urls) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/jobs`, {
        url: urls.filter((url) => url.trim() !== ""),
      });
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.message || "Failed to submit URLs");
    }
  },

  searchContent: async (query, category, page = 1, pageSize = 10) => {
    try {
      const params = new URLSearchParams();
      if (query) params.append("q", query);
      if (category) params.append("category", category);
      params.append("page", page);
      params.append("page_size", pageSize);

      const response = await axios.get(`${SEARCH_BASE_URL}/search?${params}`);
      return response.data;
    } catch (error) {
      throw new Error(error.response?.data?.error || "Search failed");
    }
  },

  checkHealth: async () => {
    try {
      await axios.get(`${API_BASE_URL}/health`);
      await axios.get(`${SEARCH_BASE_URL}/health`);
      return true;
    } catch (error) {
      return false;
    }
  },
};

