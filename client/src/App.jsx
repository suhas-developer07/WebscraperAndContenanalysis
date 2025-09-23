import React, { useState } from "react";
import Navbar from "./components/Navbar";
import SearchForm from "./components/SearchForm";
import SearchResults from "./components/SearchResults";
import UrlSubmitForm from "./components/UrlSubmitForm";
import { api } from "./services/api";
import "./App.css";

function App() {
  const [activeTab, setActiveTab] = useState("search"); // "search" or "submit"
  const [searchResults, setSearchResults] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const [page, setPage] = useState(1);

  const handleSearch = async ({ query, category }) => {
    setIsLoading(true);
    setPage(1);

    try {
      const response = await api.searchContent({
        query,
        category,
        page: 1,
      });

      setSearchResults(response);
      setHasMore(response.has_more || false);
    } catch (error) {
      console.error("Search error:", error);
      setSearchResults({ results: [] });
    } finally {
      setIsLoading(false);
    }
  };

  const loadMoreResults = async () => {
    if (!searchResults) return;
    const nextPage = page + 1;

    try {
      const response = await api.searchContent({
        query: searchResults.query,
        category: searchResults.category,
        page: nextPage,
      });

      setSearchResults((prev) => ({
        ...response,
        results: [...prev.results, ...response.results],
      }));
      setPage(nextPage);
      setHasMore(response.has_more || false);
    } catch (error) {
      console.error("Load more error:", error);
    }
  };

  return (
    <div className="App">
      <Navbar activeTab={activeTab} setActiveTab={setActiveTab} />

      <div className="content">
        {activeTab === "search" && (
          <>
            <SearchForm onSearch={handleSearch} isLoading={isLoading} />
            <SearchResults
              results={searchResults}
              isLoading={isLoading}
              hasMore={hasMore}
              onLoadMore={loadMoreResults}
            />
          </>
        )}

        {activeTab === "submit" && <UrlSubmitForm />}
      </div>
    </div>
  );
}

export default App;
