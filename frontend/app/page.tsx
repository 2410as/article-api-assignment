"use client";

import Image from "next/image";
import { useEffect, useState } from "react";

interface Article {
  id: number;
  title: string;
  body: string;
  is_pinned: boolean;
  medias: {
    id: number
    content_url: string;
  }[]
}

export default function Home() {
  const [articles, setArticles] = useState<Article[]>([]);
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");

  // 記事取得
  const fetchArticles = async () => {
    try {
      const res = await fetch("http://localhost:8080/articles");
      const data = await res.json();
      setArticles(data || []);
    } catch (e) {
      console.error("Failed to fetch articles", e);
    }
  };

  useEffect(() => {
    fetchArticles();
  }, []);

  // 記事保存
  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await fetch("http://localhost:8080/articles", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, body }),
      });
      if (!res.ok) throw new Error(`HTTP error! status: ${res.status}`);
      setTitle("");
      setBody("");
      fetchArticles();
    } catch (err) {
      console.error("記事の保存に失敗しました", err);
    }
  };

  const handleDelete = async (id: number, e: React.MouseEvent) => {
  e.stopPropagation();
  if (!confirm("この記事を完全に削除しますか？")) return;

  const res = await fetch(`http://localhost:8080/articles/${id}`, { method: "DELETE" });
  if (!res.ok) {
    console.error("削除失敗", await res.text());
    return;
  }

  // 追加：削除した記事が左側で選択中なら表示を消す
  setSelectedArticle(prev => (prev?.id === id ? null : prev));

  fetchArticles();
};


  // ピン留め切り替え
  const handleTogglePin = async (id: number, e: React.MouseEvent) => {
    e.stopPropagation();
    await fetch(`http://localhost:8080/articles/${id}/pin`, { method: "PATCH" });
    fetchArticles();
  };

  const pinnedArticles = articles.filter(a => a.is_pinned);

  return (
    <main className="flex h-screen bg-white">
      {/* 左側 */}
      <div className="w-1/2 p-10 border-r overflow-y-auto">
        {/* 記事作成フォーム */}
        <form onSubmit={handleSave} className="space-y-4 pb-10 border-b">
          <input
            type="text"
            placeholder="タイトル"
            className="w-full p-3 border rounded-lg outline-none focus:ring-2 focus:ring-green-400"
            value={title}
            onChange={e => setTitle(e.target.value)}
          />
          <textarea
            placeholder="本文を入力してください..."
            rows={6}
            className="w-full p-3 border rounded-lg outline-none focus:ring-2 focus:ring-green-400"
            value={body}
            onChange={e => setBody(e.target.value)}
          />
          <button className="w-full bg-green-500 text-white font-bold py-3 rounded-lg hover:bg-green-600 transition">
            保存する
          </button>
        </form>

        {/* ピン留め記事 */}
        <div className="mt-8">
          <h3 className="text-sm font-bold text-gray-400 mb-4 uppercase tracking-wider">
            ピン留めした記事
          </h3>
          <div className="flex flex-wrap gap-2">
            {pinnedArticles.length > 0 ? (
              pinnedArticles.map(art => (
                <div
                  key={art.id}
                  onClick={() => setSelectedArticle(art)}
                  className="px-4 py-2 bg-green-50 border border-green-200 rounded-full text-green-700 text-sm font-medium cursor-pointer hover:bg-green-100 transition"
                >
                  📌 {art.title}
                </div>
              ))
            ) : (
              <p className="text-gray-300 text-sm italic">
                ピン留めされた記事はありません
              </p>
            )}
          </div>
        </div>

        {/* 選択中の記事 */}
        <div className="mt-12 pt-8 border-t border-gray-50">
          {selectedArticle ? (
            <div>
              <h1 className="text-3xl font-black text-gray-900 mb-6">
                {selectedArticle.title}
              </h1>
              <div className="text-gray-600 leading-relaxed whitespace-pre-wrap">
                {selectedArticle.body}
              </div>
              <div className="flex flex-wrap gap-2 mt-6">
                {selectedArticle.medias.map(m => {
                  return (<Image key={m.id} width={200} height={200} alt="Picture of media" src={m.content_url} />)
                })}
              </div>
            </div>
          ) : (
            <div className="py-20 text-center text-gray-300 border-2 border-dashed rounded-xl">
              記事を選択して表示
            </div>
          )}
        </div>
      </div>

      {/* 右側：全記事リスト */}
      <div className="w-1/2 p-10 bg-gray-50 overflow-y-auto">
        <h2 className="text-xl font-bold mb-6 text-gray-700">すべての記事</h2>
        <div className="space-y-3">
          {articles.map(art => (
            <div
              key={art.id}
              onClick={() => setSelectedArticle(art)}
              className="group p-5 bg-white rounded-xl border border-gray-100 shadow-sm hover:border-green-300 transition-all relative cursor-pointer"
            >
              <h3 className="font-bold text-gray-800 pr-16">{art.title}</h3>
              <div className="absolute top-4 right-4 flex space-x-2">
                <button
                  onClick={e => handleTogglePin(art.id, e)}
                  className={`text-lg transition ${
                    art.is_pinned
                      ? "opacity-100"
                      : "opacity-0 group-hover:opacity-50"
                  }`}
                >
                  {art.is_pinned ? "📌" : "📍"}
                </button>
                <button
                  onClick={e => handleDelete(art.id, e)}
                  className="text-xs bg-red-50 text-red-400 px-2 py-1 rounded opacity-0 group-hover:opacity-100 hover:bg-red-100"
                >
                  削除
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
