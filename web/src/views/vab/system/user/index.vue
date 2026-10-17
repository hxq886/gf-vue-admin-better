<template>
  <div class="system-user-container">
    <el-card shadow="never">
      <!-- 搜索栏 -->
      <el-form :inline="true" :model="queryForm" @submit.prevent>
        <el-form-item label="账号">
          <el-input
            v-model="queryForm.username"
            placeholder="请输入账号模糊查询"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="queryForm.status"
            placeholder="全部"
            clearable
            style="width: 120px"
          >
            <el-option label="启用" :value="1" />
            <el-option label="停用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 工具栏 -->
      <div class="table-toolbar">
        <el-button
          v-if="checkPermission(['system:user:create'])"
          type="primary"
          @click="handleCreate"
        >
          新增用户
        </el-button>
      </div>

      <!-- 数据表格 -->
      <el-table v-loading="loading" :data="list" border stripe>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="username" label="账号" min-width="110" />
        <el-table-column prop="nickname" label="昵称" min-width="110" />
        <el-table-column label="角色" min-width="140">
          <template #default="{ row }">
            <el-tag
              v-for="roleId in row.roleIds"
              :key="roleId"
              size="small"
              style="margin-right: 5px"
            >
              {{ roleMap[roleId] || roleId }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="150" />
        <el-table-column prop="mobile" label="手机号" min-width="110" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? "启用" : "停用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="createdAt"
          label="创建时间"
          width="160"
          align="center"
        />
        <el-table-column label="操作" width="200" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="checkPermission(['system:user:update'])"
              type="primary"
              link
              @click="handleEdit(row)"
            >
              编辑
            </el-button>
            <el-button
              v-if="checkPermission(['system:user:password'])"
              type="warning"
              link
              @click="handlePassword(row)"
            >
              重置密码
            </el-button>
            <el-button
              v-if="checkPermission(['system:user:delete'])"
              type="danger"
              link
              @click="handleDelete(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="table-pagination"
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="total"
        :page-sizes="[10, 20, 50]"
        v-model:current-page="queryForm.page"
        v-model:page-size="queryForm.pageSize"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑用户' : '新增用户'"
      width="560px"
      destroy-on-close
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item v-if="!form.id" label="账号" prop="username">
          <el-input v-model="form.username" placeholder="请输入登录账号" />
        </el-form-item>
        <el-form-item v-if="!form.id" label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="请输入登录密码"
          />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" placeholder="请输入用户昵称" />
        </el-form-item>
        <el-form-item label="角色" prop="roleIds">
          <el-select v-model="form.roleIds" multiple placeholder="请选择角色" style="width: 100%">
            <el-option
              v-for="role in roleList"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="form.mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
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

    <!-- 重置密码对话框 -->
    <el-dialog v-model="passwordVisible" title="重置密码" width="420px" destroy-on-close>
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="90px">
        <el-form-item label="账号">
          <el-input :model-value="passwordForm.username" disabled />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input
            v-model="passwordForm.password"
            type="password"
            show-password
            placeholder="请输入新密码(6-32位)"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="passwordVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handlePasswordSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
  import { ElMessage, ElMessageBox } from "element-plus";
  import checkPermission from "@/utils/permission";
  import {
    getUserList,
    createUser,
    updateUser,
    resetUserPassword,
    deleteUser,
    getRoleList,
  } from "@/api/system";

  export default {
    name: "SystemUser",
    data() {
      return {
        loading: false,
        submitting: false,
        list: [],
        total: 0,
        roleList: [],
        queryForm: { page: 1, pageSize: 10, username: "", status: undefined },
        dialogVisible: false,
        form: this.getDefaultForm(),
        rules: {
          username: [
            { required: true, message: "请输入登录账号", trigger: "blur" },
            { min: 2, max: 64, message: "账号长度为2-64位", trigger: "blur" },
          ],
          password: [
            { required: true, message: "请输入登录密码", trigger: "blur" },
            { min: 6, max: 32, message: "密码长度为6-32位", trigger: "blur" },
          ],
          nickname: [{ required: true, message: "请输入昵称", trigger: "blur" }],
          email: [{ type: "email", message: "邮箱格式不正确", trigger: "blur" }],
        },
        passwordVisible: false,
        passwordForm: { id: undefined, username: "", password: "" },
        passwordRules: {
          password: [
            { required: true, message: "请输入新密码", trigger: "blur" },
            { min: 6, max: 32, message: "密码长度为6-32位", trigger: "blur" },
          ],
        },
      };
    },
    computed: {
      // 角色ID -> 角色名称映射
      roleMap() {
        const map = {};
        this.roleList.forEach((role) => (map[role.id] = role.name));
        return map;
      },
    },
    created() {
      this.fetchList();
      this.fetchRoles();
    },
    methods: {
      checkPermission,
      getDefaultForm() {
        return {
          id: undefined,
          username: "",
          password: "",
          nickname: "",
          email: "",
          mobile: "",
          status: 1,
          remark: "",
          roleIds: [],
        };
      },
      // 查询用户列表
      async fetchList() {
        this.loading = true;
        try {
          const { data } = await getUserList(this.queryForm);
          this.list = data.list || [];
          this.total = data.total || 0;
        } finally {
          this.loading = false;
        }
      },
      // 查询角色列表 (供下拉与标签展示)
      async fetchRoles() {
        const { data } = await getRoleList();
        this.roleList = data.list || [];
      },
      handleSearch() {
        this.queryForm.page = 1;
        this.fetchList();
      },
      handleReset() {
        this.queryForm = { page: 1, pageSize: 10, username: "", status: undefined };
        this.fetchList();
      },
      handleCreate() {
        this.form = this.getDefaultForm();
        this.dialogVisible = true;
      },
      handleEdit(row) {
        this.form = {
          id: row.id,
          username: row.username,
          nickname: row.nickname,
          email: row.email,
          mobile: row.mobile,
          status: row.status,
          remark: row.remark,
          roleIds: (row.roleIds || []).slice(),
        };
        this.dialogVisible = true;
      },
      // 提交新增/编辑
      handleSubmit() {
        this.$refs.formRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            if (this.form.id) {
              const { password, username, ...data } = this.form;
              await updateUser(data);
              ElMessage.success("修改成功");
            } else {
              await createUser(this.form);
              ElMessage.success("新增成功");
            }
            this.dialogVisible = false;
            this.fetchList();
          } finally {
            this.submitting = false;
          }
        });
      },
      handlePassword(row) {
        this.passwordForm = { id: row.id, username: row.username, password: "" };
        this.passwordVisible = true;
      },
      handlePasswordSubmit() {
        this.$refs.passwordFormRef.validate(async (valid) => {
          if (!valid) return;
          this.submitting = true;
          try {
            await resetUserPassword({
              id: this.passwordForm.id,
              password: this.passwordForm.password,
            });
            ElMessage.success("密码重置成功");
            this.passwordVisible = false;
          } finally {
            this.submitting = false;
          }
        });
      },
      // 删除用户 (内置超管受后端保护)
      handleDelete(row) {
        ElMessageBox.confirm(`确认删除用户 [${row.username}] 吗?`, "提示", {
          type: "warning",
        }).then(async () => {
          await deleteUser({ ids: [row.id] });
          ElMessage.success("删除成功");
          this.fetchList();
        });
      },
    },
  };
</script>

<style lang="scss" scoped>
  .system-user-container {
    padding: 15px;

    .table-toolbar {
      margin-bottom: 15px;
    }

    .table-pagination {
      margin-top: 15px;
      justify-content: flex-end;
    }
  }
</style>
