const { createApp, ref, computed, onMounted, watch } = Vue;

createApp({
  setup() {
    // 响应式数据
    const allImages = ref([]);
    const filteredImages = ref([]);
    const searchFilters = ref({ name: "", tag: "", desc: "" });
    const loading = ref(false);
    const errorMsg = ref("");
    const showModal = ref(false);
    const modalImage = ref(null);
    const searchDebounce = ref(null);

    // 加载图片数据
    async function loadImageData() {
      loading.value = true;
      errorMsg.value = "";

      try {
        // 获取索引文件
        const indexResponse = await fetch(`${window.location.origin}/records/index.json`);

        if (!indexResponse.ok) {
          throw new Error(`HTTP error! status: ${indexResponse.status}`);
        }

        const files = await indexResponse.json();
        const allData = [];

        // 加载所有记录文件
        for (const file of files) {
          try {
            const res = await fetch(`${window.location.origin}/records/${file}`);
            if (res.ok) {
              const data = await res.json();
              allData.push(...data);
            } else {
              console.warn(`无法加载文件: ${file}, status: ${res.status}`);
            }
          } catch (e) {
            console.warn(`加载文件失败: ${file}`, e);
          }
        }

        // 按创建时间排序（最新的在前）
        allImages.value = allData.sort((a, b) => b.created_at - a.created_at);
        filteredImages.value = [...allImages.value];

      } catch (e) {
        console.error("加载数据失败:", e);
        errorMsg.value = "加载数据失败，请检查网络连接或刷新页面重试";
        loadMockData();
      } finally {
        loading.value = false;
      }
    }

    // 模拟数据（当真实数据加载失败时使用）
    function loadMockData() {
      allImages.value = [
        {
          filename: "示例图片.jpg",
          url: "https://via.placeholder.com/800x600/3498db/ffffff?text=示例图片",
          thumb_url: "https://via.placeholder.com/300x200/3498db/ffffff?text=示例",
          created_at: Date.now() / 1000,
          tags: ["示例", "测试"],
          desc: "这是一个示例图片，用于演示功能",
          size_kb: 1024,
          width: 800,
          height: 600,
          hash: "mock_hash_123",
        },
      ];
      filteredImages.value = [...allImages.value];
    }

    // 执行搜索（带防抖）
    function performSearch() {
      if (searchDebounce.value) {
        clearTimeout(searchDebounce.value);
      }

      searchDebounce.value = setTimeout(() => {
        filteredImages.value = allImages.value.filter((img) => {
          const nameMatch = !searchFilters.value.name ||
            img.filename.toLowerCase().includes(searchFilters.value.name.toLowerCase());

          const tagMatch = !searchFilters.value.tag ||
            img.tags.some((t) => t.toLowerCase().includes(searchFilters.value.tag.toLowerCase()));

          const descMatch = !searchFilters.value.desc ||
            (img.desc && img.desc.toLowerCase().includes(searchFilters.value.desc.toLowerCase()));

          return nameMatch && tagMatch && descMatch;
        });
      }, 300);
    }

    // 清空搜索
    function clearSearch() {
      searchFilters.value = { name: "", tag: "", desc: "" };
      filteredImages.value = [...allImages.value];
    }

    // 复制到剪贴板
    async function copyToClipboard(text) {
      try {
        await navigator.clipboard.writeText(text);
        showToast("已复制链接到剪贴板！");
      } catch (err) {
        console.error("复制失败:", err);
        showToast("复制失败，请手动复制");
      }
    }

    // 显示提示信息
    function showToast(message) {
      const toast = document.createElement('div');
      toast.className = 'toast';
      toast.textContent = message;
      document.body.appendChild(toast);

      setTimeout(() => {
        toast.classList.add('show');
      }, 100);

      setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
          document.body.removeChild(toast);
        }, 300);
      }, 2000);
    }

    // 打开图片详情模态框
    function openModal(img) {
      modalImage.value = img;
      showModal.value = true;
      document.body.style.overflow = 'hidden';
    }

    // 关闭模态框
    function closeModal() {
      showModal.value = false;
      document.body.style.overflow = 'auto';
    }

    // 格式化文件大小
    function formatFileSize(bytes) {
      if (bytes < 1024) return bytes + ' B';
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
      return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
    }

    // 格式化日期
    function formatDate(timestamp) {
      return new Date(timestamp * 1000).toLocaleDateString("zh-CN", {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    }

    // 计算属性
    const totalCount = computed(() => allImages.value.length);
    const filteredCount = computed(() => filteredImages.value.length);
    const totalSize = computed(() => {
      const totalBytes = allImages.value.reduce((sum, img) => sum + (img.size_kb * 1024), 0);
      return formatFileSize(totalBytes);
    });

    // 监听搜索条件变化
    watch(searchFilters, performSearch, { deep: true });

    // 组件挂载时加载数据
    onMounted(loadImageData);

    return {
      allImages,
      filteredImages,
      searchFilters,
      loading,
      errorMsg,
      performSearch,
      clearSearch,
      totalCount,
      filteredCount,
      totalSize,
      showModal,
      modalImage,
      openModal,
      closeModal,
      copyToClipboard,
      formatFileSize,
      formatDate,
      loadImageData,
    };
  },

  template: `
    <div class="container">
      <!-- 页眉 -->
      <header class="header">
        <h1>📸 图床管理工具</h1>
        <p>图片上传记录展示与管理</p>
      </header>

      <!-- 错误提示 -->
      <div v-if="errorMsg" class="error-message">
        <span>⚠️ {{ errorMsg }}</span>
        <button @click="loadImageData">重试</button>
      </div>

      <!-- 搜索区域 -->
      <section class="search-section">
        <div class="search-form">
          <div class="search-group">
            <label>文件名搜索:</label>
            <input 
              v-model="searchFilters.name" 
              placeholder="输入文件名关键词..."
              type="text"
            >
          </div>
          <div class="search-group">
            <label>标签搜索:</label>
            <input 
              v-model="searchFilters.tag" 
              placeholder="输入标签关键词..."
              type="text"
            >
          </div>
          <div class="search-group">
            <label>描述搜索:</label>
            <input 
              v-model="searchFilters.desc" 
              placeholder="输入描述关键词..."
              type="text"
            >
          </div>
          <div class="search-buttons">
            <button @click="performSearch" class="btn-primary">🔍 搜索</button>
            <button @click="clearSearch" class="btn-secondary">🗑️ 清空</button>
          </div>
        </div>
      </section>

      <!-- 统计信息 -->
      <section class="stats-section">
        <div class="stats">
          <div class="stat-item">
            <span class="stat-label">总图片数:</span>
            <span class="stat-value">{{ totalCount }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">显示图片数:</span>
            <span class="stat-value">{{ filteredCount }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">总大小:</span>
            <span class="stat-value">{{ totalSize }}</span>
          </div>
        </div>
      </section>

      <!-- 图片展示 -->
      <section class="gallery-section">
        <div v-if="loading" class="loading">
          <div class="loading-spinner"></div>
          <p>正在加载数据...</p>
        </div>
        
        <div v-else-if="filteredImages.length === 0" class="no-results">
          <div class="no-results-icon">📷</div>
          <p>没有找到匹配的图片</p>
          <button @click="clearSearch" class="btn-primary">查看所有图片</button>
        </div>
        
        <div v-else class="gallery">
          <div 
            v-for="img in filteredImages" 
            :key="img.hash" 
            class="image-card" 
            @click="openModal(img)"
          >
            <div class="image-thumbnail-container">
              <img 
                :src="img.thumb_url" 
                :alt="img.filename" 
                class="image-thumbnail"
                @error="$event.target.src='https://via.placeholder.com/300x200/eee/999?text=加载失败'"
              >
            </div>
            <div class="image-info">
              <div class="image-title" :title="img.filename">{{ img.filename }}</div>
              <div class="image-meta">
                <div>📅 {{ formatDate(img.created_at) }}</div>
                <div>📏 {{ img.width }} × {{ img.height }}</div>
                <div>💾 {{ formatFileSize(img.size_kb * 1024) }}</div>
                <div v-if="img.desc" class="image-desc">📝 {{ img.desc }}</div>
              </div>
              <div v-if="img.tags && img.tags.length" class="image-tags">
                <span v-for="tag in img.tags" :key="tag" class="tag">
                  {{ tag }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 图片详情模态框 -->
      <div v-if="showModal" class="modal" @click.self="closeModal">
        <div class="modal-content">
          <span class="close" @click="closeModal">&times;</span>
          <div v-if="modalImage" class="modal-body">
            <div class="modal-image-container">
              <img 
                :src="modalImage.url" 
                :alt="modalImage.filename" 
                class="modal-image"
                @error="e => e.target.src = 'https://via.placeholder.com/600x400/eee/999?text=图片加载失败'"
              >
            </div>
            <div class="modal-info">
              <div class="info-row">
                <span class="info-label">文件名:</span>
                <span class="info-value">{{ modalImage.filename }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">创建时间:</span>
                <span class="info-value">{{ formatDate(modalImage.created_at) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">尺寸:</span>
                <span class="info-value">{{ modalImage.width }} × {{ modalImage.height }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">文件大小:</span>
                <span class="info-value">{{ formatFileSize(modalImage.size_kb * 1024) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">图片链接:</span>
                <div class="info-value">
                  <a :href="modalImage.url" target="_blank" class="link">{{ modalImage.url }}</a>
                  <button @click.stop="copyToClipboard(modalImage.url)" class="btn-copy">复制</button>
                </div>
              </div>
              <div v-if="modalImage.desc" class="info-row">
                <span class="info-label">描述:</span>
                <span class="info-value">{{ modalImage.desc }}</span>
              </div>
              <div v-if="modalImage.tags && modalImage.tags.length" class="info-row">
                <span class="info-label">标签:</span>
                <div class="info-value">
                  <span v-for="tag in modalImage.tags" :key="tag" class="tag">{{ tag }}</span>
                </div>
              </div>
              <div class="info-row">
                <span class="info-label">哈希值:</span>
                <span class="info-value hash">{{ modalImage.hash }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
}).mount("#app");
