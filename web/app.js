/**
 * 图床管理工具前端展示页面
 * 功能：读取records目录下的JSON文件，展示图片信息，支持搜索功能
 */

class PhotoGallery {
  constructor() {
    this.allImages = [];
    this.filteredImages = [];
    this.searchFilters = {
      name: "",
      tag: "",
      desc: "",
    };
    this.init();
  }

  /**
   * 初始化应用
   */
  async init() {
    try {
      await this.loadImageData();
      this.setupEventListeners();
      this.renderGallery();
      this.updateStats();
    } catch (error) {
      console.error("初始化失败:", error);
      this.showError("加载数据失败，请检查网络连接");
    }
  }

  /**
   * 加载图片数据
   */
  async loadImageData() {
    const loadingElement = document.getElementById("loading");
    if (loadingElement) {
      loadingElement.style.display = "block";
    }

    try {
      const jsonFiles = ["2025-08.json"];
      const allData = [];

      for (const file of jsonFiles) {
        try {
          const fileResponse = await fetch(`/records/${file}`);
          if (fileResponse.ok) {
            const data = await fileResponse.json();
            allData.push(...data);
          }
        } catch (error) {
          console.warn(`无法加载文件 ${file}:`, error);
        }
      }

      this.allImages = allData;
      this.filteredImages = [...this.allImages];
    } catch (error) {
      console.error("加载图片数据失败:", error);
      this.loadMockData();
    } finally {
      if (loadingElement) {
        loadingElement.style.display = "none";
      }
    }
  }

  /**
   * 加载模拟数据
   */
  loadMockData() {
    this.allImages = [
      {
        filename: "test2.b9de6437.jpg",
        url: "https://cdn.jsdelivr.net/gh/LucienVen/photo-bed@main/images/test2.b9de6437.jpg",
        thumb_url:
          "https://cdn.jsdelivr.net/gh/LucienVen/photo-bed@main/images/test2.b9de6437.thumb.jpg",
        created_at: 1754620115,
        tags: [""],
        desc: "",
        size_kb: 5174,
        width: 4000,
        height: 3000,
        hash: "b9de6437b06a74a2d237939c8706b6d41dd53212aefa75e17891ff85f848f4de",
      },
      {
        filename: "0mCHGeTPqzw.c83a0ddc.jpg",
        url: "https://cdn.jsdelivr.net/gh/LucienVen/photo-bed@main/images/0mCHGeTPqzw.c83a0ddc.jpg",
        thumb_url:
          "https://cdn.jsdelivr.net/gh/LucienVen/photo-bed@main/images/0mCHGeTPqzw.c83a0ddc.thumb.jpg",
        created_at: 1754621331,
        tags: ["测试", "随机"],
        desc: "测试图片",
        size_kb: 2848,
        width: 4000,
        height: 6000,
        hash: "c83a0ddc2f9adc5caf120c3f0f8b9f63a63bb4adf7ad07860e36cafd0a957bcb",
      },
    ];
    this.filteredImages = [...this.allImages];
  }

  /**
   * 设置事件监听器
   */
  setupEventListeners() {
    const searchBtn = document.getElementById("searchBtn");
    if (searchBtn) {
      searchBtn.addEventListener("click", () => this.performSearch());
    }

    const clearBtn = document.getElementById("clearBtn");
    if (clearBtn) {
      clearBtn.addEventListener("click", () => this.clearSearch());
    }

    const searchInputs = ["nameSearch", "tagSearch", "descSearch"];
    searchInputs.forEach((id) => {
      const input = document.getElementById(id);
      if (input) {
        input.addEventListener("keypress", (e) => {
          if (e.key === "Enter") {
            this.performSearch();
          }
        });
      }
    });

    const modal = document.getElementById("imageModal");
    const closeBtn = document.querySelector(".close");
    if (modal && closeBtn) {
      closeBtn.addEventListener("click", () => this.closeModal());
      modal.addEventListener("click", (e) => {
        if (e.target === modal) {
          this.closeModal();
        }
      });
    }
  }

  /**
   * 执行搜索
   */
  performSearch() {
    const nameInput = document.getElementById("nameSearch");
    const tagInput = document.getElementById("tagSearch");
    const descInput = document.getElementById("descSearch");

    this.searchFilters = {
      name: nameInput?.value.toLowerCase() || "",
      tag: tagInput?.value.toLowerCase() || "",
      desc: descInput?.value.toLowerCase() || "",
    };

    this.filterImages();
    this.renderGallery();
    this.updateStats();
  }

  /**
   * 清空搜索条件
   */
  clearSearch() {
    const nameInput = document.getElementById("nameSearch");
    const tagInput = document.getElementById("tagSearch");
    const descInput = document.getElementById("descSearch");

    if (nameInput) nameInput.value = "";
    if (tagInput) tagInput.value = "";
    if (descInput) descInput.value = "";

    this.searchFilters = { name: "", tag: "", desc: "" };
    this.filteredImages = [...this.allImages];
    this.renderGallery();
    this.updateStats();
  }

  /**
   * 过滤图片
   */
  filterImages() {
    this.filteredImages = this.allImages.filter((image) => {
      const nameMatch =
        !this.searchFilters.name ||
        image.filename.toLowerCase().includes(this.searchFilters.name);

      const tagMatch =
        !this.searchFilters.tag ||
        image.tags.some((tag) =>
          tag.toLowerCase().includes(this.searchFilters.tag)
        );

      const descMatch =
        !this.searchFilters.desc ||
        image.desc.toLowerCase().includes(this.searchFilters.desc);

      return nameMatch && tagMatch && descMatch;
    });
  }

  /**
   * HTML转义函数
   */
  escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  /**
   * 渲染图片画廊
   */
  renderGallery() {
    const gallery = document.getElementById("gallery");
    const noResults = document.getElementById("noResults");

    if (!gallery) return;

    if (this.filteredImages.length === 0) {
      gallery.innerHTML = "";
      if (noResults) {
        noResults.style.display = "block";
      }
      return;
    }

    if (noResults) {
      noResults.style.display = "none";
    }

    gallery.innerHTML = this.filteredImages
      .map((image) => this.createImageCard(image))
      .join("");

    const imageCards = gallery.querySelectorAll(".image-card");
    imageCards.forEach((card, index) => {
      card.addEventListener("click", () =>
        this.showImageDetail(this.filteredImages[index])
      );
    });
  }

  /**
   * 创建图片卡片HTML
   */
  createImageCard(image) {
    const tags = image.tags
      .filter((tag) => tag.trim() !== "")
      .map((tag) => `<span class="tag">${this.escapeHtml(tag)}</span>`)
      .join("");

    const createdDate = new Date(image.created_at * 1000).toLocaleDateString(
      "zh-CN"
    );
    const sizeMB = (image.size_kb / 1024).toFixed(2);

    return `
            <div class="image-card">
                <img src="${image.thumb_url}" alt="${this.escapeHtml(
      image.filename
    )}" 
                     class="image-thumbnail" onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZjhmOWZhIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxNCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPua1j+iniOWZqDwvdGV4dD48L3N2Zz4='">
                <div class="image-info">
                    <div class="image-title">${this.escapeHtml(
                      image.filename
                    )}</div>
                    <div class="image-meta">
                        <div>📅 ${createdDate}</div>
                        <div>📏 ${image.width} × ${image.height}</div>
                        <div>💾 ${sizeMB} MB</div>
                        ${
                          image.desc
                            ? `<div>📝 ${this.escapeHtml(image.desc)}</div>`
                            : ""
                        }
                    </div>
                    ${tags ? `<div class="image-tags">${tags}</div>` : ""}
                </div>
            </div>
        `;
  }

  /**
   * 显示图片详情模态框
   */
  showImageDetail(image) {
    const modal = document.getElementById("imageModal");
    const modalContent = document.getElementById("modalContent");

    if (!modal || !modalContent) return;

    const createdDate = new Date(image.created_at * 1000).toLocaleString(
      "zh-CN"
    );
    const sizeMB = (image.size_kb / 1024).toFixed(2);
    const tags = image.tags
      .filter((tag) => tag.trim() !== "")
      .map((tag) => `<span class="tag">${this.escapeHtml(tag)}</span>`)
      .join("");

    modalContent.innerHTML = `
            <img src="${image.url}" alt="${this.escapeHtml(
      image.filename
    )}" class="modal-image" 
                 onerror="this.src='data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAwIiBoZWlnaHQ9IjQwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZjhmOWZhIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxOCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPua1j+iniOWZqDwvdGV4dD48L3N2Zz4='">
            <div class="modal-info">
                <div class="modal-info-item">
                    <strong>文件名</strong>
                    ${this.escapeHtml(image.filename)}
                </div>
                <div class="modal-info-item">
                    <strong>创建时间</strong>
                    ${createdDate}
                </div>
                <div class="modal-info-item">
                    <strong>尺寸</strong>
                    ${image.width} × ${image.height}
                </div>
                <div class="modal-info-item">
                    <strong>文件大小</strong>
                    ${sizeMB} MB
                </div>
                <div class="modal-info-item">
                    <strong>图片链接</strong>
                    <a href="${
                      image.url
                    }" target="_blank" style="word-break: break-all; color: #667eea;">${
      image.url
    }</a>
                </div>
                ${
                  image.desc
                    ? `
                <div class="modal-info-item">
                    <strong>描述</strong>
                    ${this.escapeHtml(image.desc)}
                </div>
                `
                    : ""
                }
                ${
                  tags
                    ? `
                <div class="modal-info-item">
                    <strong>标签</strong>
                    <div style="margin-top: 8px;">${tags}</div>
                </div>
                `
                    : ""
                }
                <div class="modal-info-item">
                    <strong>哈希值</strong>
                    <span style="word-break: break-all; font-family: monospace; font-size: 12px;">${
                      image.hash
                    }</span>
                </div>
            </div>
        `;

    modal.style.display = "block";
  }

  /**
   * 关闭模态框
   */
  closeModal() {
    const modal = document.getElementById("imageModal");
    if (modal) {
      modal.style.display = "none";
    }
  }

  /**
   * 更新统计信息
   */
  updateStats() {
    const totalCount = document.getElementById("totalCount");
    const filteredCount = document.getElementById("filteredCount");
    const totalSize = document.getElementById("totalSize");

    if (totalCount) {
      totalCount.textContent = this.allImages.length.toString();
    }

    if (filteredCount) {
      filteredCount.textContent = this.filteredImages.length.toString();
    }

    if (totalSize) {
      const totalSizeKB = this.allImages.reduce(
        (sum, image) => sum + image.size_kb,
        0
      );
      const totalSizeMB = (totalSizeKB / 1024).toFixed(2);
      totalSize.textContent = `${totalSizeMB} MB`;
    }
  }

  /**
   * 显示错误信息
   */
  showError(message) {
    const gallery = document.getElementById("gallery");
    if (gallery) {
      gallery.innerHTML = `
                <div style="text-align: center; padding: 40px; color: #e74c3c;">
                    <h3>❌ 错误</h3>
                    <p>${message}</p>
                </div>
            `;
    }
  }
}

// 页面加载完成后初始化应用
document.addEventListener("DOMContentLoaded", () => {
  new PhotoGallery();
});
