<template>
  <div class="system-menu-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <div class="table-toolbar">
        <el-button
          v-if="checkPermission(['system:menu:create'])"
          type="primary"
          @click="handleCreate(0)"
        >
          新增顶级菜单
        </el-button>
        <el-button @click="toggleExpand">
          {{ isExpand ? "折叠全部" : "展开全部" }}
        </el-button>
      </div>

      <!-- 数据表格 (树形) -->
      <el-table
        v-loading="loading"
        :key="isExpand ? 'expand' : 'collapse'"
        :data="menuTree"
        border
        row-key="id"
        :default-expand-all="isExpand"
      >
        <el-table-column prop="title" label="菜单标题" min-width="160">
          <template #default="{ row }">
            <el-tag v-if="row.menuType === 'M'" size="small">目录</el-tag>
            <el-tag v-else-if="row.menuType === 'C'" type="success" size="small">
              菜单
            </el-tag>
            <el-tag v-else type="info" size="small">按钮</el-tag>
            <span style="margin-left: 8px">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="icon" label="图标" width="90" align="center">
          <template #default="{ row }">
            <span v-if="row.icon">{{ row.icon }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="menuType" label="类型" width="70" align="center" />
        <el-table-column prop="permission" label="权限标识" min-width="150">
          <template #default="{ row }">
            <el-tag v-if="row.permission" size="small" type="warning">
              {{ row.permission }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由path" min-width="120" />
        <el-table-column prop="component" label="组件" min-width="150" />
        <el-table-column prop="sort" label="排序" width="70" align="center" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? "启用" : "停用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="checkPermission(['system:menu:create']) && row.menuType !== 'F'"
              type="primary"
              link
              @click="handleCreate(row.id)"
            >
              新增
            </el-button>
            <el-button
              v-if="checkPermission(['system:menu:update'])"
              type="primary"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="checkPermission(['system:menu:delete'])"
              type="danger"
              link
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑菜单' : '新增菜单'"
      width="640px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="上级菜单">
          <el-tree-select
            v-model="form.parentId"
            :data="parentOptions"
            node-key="id"
            :props="{ label: 'title', children: 'children' }"
            check-strictly
            default-expand-all
            placeholder="顶级菜单"
            clearable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="菜单类型">
          <el-radio-group v-model="form.menuType">
            <el-radio value="M">目录</el-radio>
            <el-radio value="C">菜单</el-radio>
            <el-radio value="F">按钮</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入菜单标题" />
        </el-form-item>
        <template v-if="form.menuType !== 'F'">
          <el-form-item label="路由name">
            <el-input v-model="form.name" placeholder="请输入路由name(唯一, 如 SystemUser)" />
          </el-form-item>
          <el-form-item label="路由path">
            <el-input v-model="form.path" placeholder="请输入路由path(如 system/user)" />
          </el-form-item>
          <el-form-item v-if="form.menuType === 'C'" label="组件路径">
            <el-input
              v-model="form.component"
              placeholder="视图相对路径(如 vab/system/user/index)"
            />
          </el-form-item>
          <el-form-item label="图标">
            <el-input v-model="form.icon" placeholder="请输入图标名(FA风格, 如 user)" />
          </el-form-item>
        </template>
        <el-form-item v-if="form.menuType === 'F'" label="权限标识">
          <el-input
            v-model="form.permission"
            placeholder="请输入权限标识(如 system:user:create)"
          />
        </el-form-item>
        <el-form-item label="显示顺序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.menuType === 'C'" label="页面缓存">
          <el-switch v-model="form.keepAlive" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
  import { ElMessage, ElMessageBox } from "element-plus";
  import checkPermission from "@/utils/permission";
  import { getMenuList, createMenu, updateMenu, deleteMenu } from "@/api/system";

  export default {
    name: "SystemMenu",
    data() {
      return {
        loading: false,
        submitting: false,
        list: [],
        isExpand: true,
        dialogVisible: false,
        form: this.getDefaultForm(),
        rules: {
          title: [{ required: true, message: "请输入菜单标题", trigger: "blur" }],
          menuType: [{ required: true, message: "请选择菜单类型", trigger: "change" }],
        },
      };
    },
    computed: {
      // 菜单平铺列表构建树 (仅目录/菜单可作上级)
      menuTree() {
        const childrenMap = {};
        const roots = [];
        this.list.forEach((menu) => {
          menu.children = childrenMap[menu.id] || [];
          childrenMap[menu.id] = menu.children;
          if (menu.parentId === 0) roots.push(menu);
          else {
            childrenMap[menu.parentId] = childrenMap[menu.parentId] || [];
            childrenMap[menu.parentId].push(menu);
          }
        });
        return roots;
      },
      // 上级菜单选项 (附"顶级"节点, 排除按钮)
      parentOptions() {
        const top = [{ id: 0, title: "顶级菜单", children: this.menuTree }];
        return top;
      },
    },
    created() {
      this.fetchList();
    },
    methods: {
      checkPermission,
      getDefaultForm() {
        return {
          id: undefined,
          parentId: 0,
          menuType: "C",
          title: "",
          name: "",
          path: "",
          component: "",
          icon: "",
          permission: "",
          sort: 1,
          status: 1,
          keepAlive: 0,
          remark: "",
        };
      },
      async fetchList() {
        this.loading = true;
        try {
          const { data } = await getMenuList();
          this.list = data.list || [];
        } finally {
          this.loading = false;
        }
      },
      toggleExpand() {
        this.isExpand = !this.isExpand;
      },
      handleCreate(parentId) {
        this.form = this.getDefaultForm();
        this.form.parentId = parentId;
        this.dialogVisible = true;
      },
      handleEdit(row) {
        this.form = {
          id: row.id,
          parentId: row.parentId,
          menuType: row.menuType,
          title: row.title,
          name: row.name || "",
          path: row.path || "",
          component: row.component || "",
          redirect: row.redirect || "",
          icon: row.icon || "",
          permission: row.permission || "",
          alwaysShow: row.alwaysShow,
          hidden: row.hidden,
          keepAlive: row.keepAlive,
          affix: row.affix,
          sort: row.sort,
          status: row.status,
          remark: row.remark || "",
        };
        this.dialogVisible = true;
      },
      handleSubmit() {
        this.$refs.formRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            if (this.form.id) {
              await updateMenu(this.form);
              ElMessage.success("修改成功");
            } else {
              await createMenu(this.form);
              ElMessage.success("新增成功");
            }
            this.dialogVisible = false;
            this.fetchList();
          } finally {
            this.submitting = false;
          }
        });
      },
      // 删除菜单 (存在子菜单时后端拒绝)
      handleDelete(row) {
        ElMessageBox.confirm(`确认删除菜单 [${row.title}] 吗?`, "提示", {
          type: "warning",
        }).then(async () => {
          await deleteMenu({ ids: [row.id] });
          ElMessage.success("删除成功");
          this.fetchList();
        });
      },
    },
  };
</script>

<style lang="scss" scoped>
  .system-menu-container {
    padding: 15px;

    .table-toolbar {
      margin-bottom: 15px;
    }
  }
</style>
