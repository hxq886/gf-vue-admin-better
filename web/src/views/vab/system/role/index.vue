<template>
  <div class="system-role-container">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <div class="table-toolbar">
        <el-button
          v-if="checkPermission(['system:role:create'])"
          type="primary"
          @click="handleCreate"
        >
          新增角色
        </el-button>
      </div>

      <!-- 数据表格 -->
      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="name" label="角色名称" min-width="120" />
        <el-table-column prop="roleKey" label="权限字符" min-width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.roleKey }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="80" align="center" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? "启用" : "停用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" />
        <el-table-column
          prop="createdAt"
          label="创建时间"
          width="160"
          align="center"
        />
        <el-table-column label="操作" width="140" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="checkPermission(['system:role:update'])"
              type="primary"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="checkPermission(['system:role:delete'])"
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
      :title="form.id ? '编辑角色' : '新增角色'"
      width="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="权限字符" prop="roleKey">
          <el-input
            v-model="form.roleKey"
            placeholder="请输入权限字符(如: common)"
            :disabled="form.id === 1"
          />
        </el-form-item>
        <el-form-item label="显示顺序">
          <el-input-number v-model="form.sort" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status" :disabled="form.id === 1">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="菜单权限">
          <el-tree
            ref="menuTreeRef"
            :data="menuTree"
            show-checkbox
            node-key="id"
            :props="{ label: 'title', children: 'children' }"
            default-expand-all
            style="width: 100%"
          />
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
  import { getRoleList, createRole, updateRole, deleteRole, getMenuList } from "@/api/system";

  export default {
    name: "SystemRole",
    data() {
      return {
        loading: false,
        submitting: false,
        list: [],
        menuList: [],
        dialogVisible: false,
        form: this.getDefaultForm(),
        rules: {
          name: [{ required: true, message: "请输入角色名称", trigger: "blur" }],
          roleKey: [{ required: true, message: "请输入权限字符", trigger: "blur" }],
        },
      };
    },
    computed: {
      // 菜单平铺列表构建树 (供权限树展示)
      menuTree() {
        const childrenMap = {};
        const roots = [];
        this.menuList.forEach((menu) => {
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
    },
    created() {
      this.fetchList();
      this.fetchMenus();
    },
    methods: {
      checkPermission,
      getDefaultForm() {
        return { id: undefined, name: "", roleKey: "", sort: 1, status: 1, remark: "", menuIds: [] };
      },
      async fetchList() {
        this.loading = true;
        try {
          const { data } = await getRoleList();
          this.list = data.list || [];
        } finally {
          this.loading = false;
        }
      },
      async fetchMenus() {
        const { data } = await getMenuList();
        this.menuList = data.list || [];
      },
      handleCreate() {
        this.form = this.getDefaultForm();
        this.dialogVisible = true;
        this.$nextTick(() => this.$refs.menuTreeRef?.setCheckedKeys([]));
      },
      handleEdit(row) {
        this.form = {
          id: row.id,
          name: row.name,
          roleKey: row.roleKey,
          sort: row.sort,
          status: row.status,
          remark: row.remark,
          menuIds: (row.menuIds || []).slice(),
        };
        this.dialogVisible = true;
        // 仅勾选叶子节点: el-tree 勾选父节点会级联勾选全部子孙,
        // 若直接回显全部menuIds会导致保存时误授权所有子菜单;
        // 父节点随叶子勾选状态自动呈半选/全选, 保存时再合并半选父节点
        this.$nextTick(() => {
          const parentIds = new Set(this.menuList.map((m) => m.parentId));
          const leafIds = (row.menuIds || []).filter((id) => !parentIds.has(id));
          this.$refs.menuTreeRef?.setCheckedKeys(leafIds);
        });
      },
      handleSubmit() {
        this.$refs.formRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            // 取勾选节点与半选父节点(上级目录随子级授权)
            const tree = this.$refs.menuTreeRef;
            const data = {
              ...this.form,
              menuIds: [
                ...tree.getCheckedKeys(),
                ...tree.getHalfCheckedKeys(),
              ].map(Number),
            };
            if (data.id) {
              await updateRole(data);
              ElMessage.success("修改成功");
            } else {
              await createRole(data);
              ElMessage.success("新增成功");
            }
            this.dialogVisible = false;
            this.fetchList();
          } finally {
            this.submitting = false;
          }
        });
      },
      // 删除角色 (内置超管角色受后端保护)
      handleDelete(row) {
        ElMessageBox.confirm(`确认删除角色 [${row.name}] 吗?`, "提示", {
          type: "warning",
        }).then(async () => {
          await deleteRole({ ids: [row.id] });
          ElMessage.success("删除成功");
          this.fetchList();
        });
      },
    },
  };
</script>

<style lang="scss" scoped>
  .system-role-container {
    padding: 15px;

    .table-toolbar {
      margin-bottom: 15px;
    }
  }
</style>
